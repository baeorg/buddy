package storage

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"runtime"
	"time"

	"github.com/baeorg/buddy/pkg/share"
	"github.com/baeorg/buddy/pkg/types"
	"github.com/bytedance/sonic"
	"github.com/sunvim/gmdbx"
	"github.com/sunvim/mq"
	"go.uber.org/atomic"
)

var (
	defaultGeometry = gmdbx.Geometry{
		SizeLower:       4 * share.MB,
		SizeNow:         4 * share.MB,
		SizeUpper:       4 * share.TB,
		GrowthStep:      16 * share.MB,
		ShrinkThreshold: 64 * share.MB,
		PageSize:        64 * share.KB,
	}
	defaultFlags = gmdbx.EnvSyncDurable |
		gmdbx.EnvWriteMap |
		gmdbx.EnvLIFOReclaim |
		gmdbx.EnvCoalesce | gmdbx.EnvNoSubDir
)

const (
	SubDBSize       = 7
	SubDBUsers      = "users"       // store user info
	SubDBPermits    = "permits"     // store user permissions
	SubDBRels       = "user:rels"   // store user relationships
	SubDBConvsUsers = "convs:users" // store conversation users
	SubDBConvsMesgs = "convs:mesgs" // store conversation messages
	SubDBMesgs      = "mesgs"       // store messages
	SubDBEnv        = "env"         // store environment variables
)

var (
	SeqmConvs string = "convs:id" // conversation id
	SeqmMesgs string = "mesgs:id" // message id
)

type DB struct {
	mesgeq     *mq.MessageQueue
	seqm       map[string]*atomic.Uint64
	env        *gmdbx.Env
	users      gmdbx.DBI
	permits    gmdbx.DBI
	rels       gmdbx.DBI
	convsUsers gmdbx.DBI
	convsMesgs gmdbx.DBI
	mesgs      gmdbx.DBI
	genv       gmdbx.DBI
}

func New(ctx context.Context, path string, mqPath string) *DB {
	db := &DB{}

	env, err := gmdbx.NewEnv()
	if err != gmdbx.ErrSuccess {
		panic(err)
	}
	db.env = env

	err = env.SetMaxDBS(SubDBSize)
	if err != gmdbx.ErrSuccess {
		log.Fatal("set max dbs failed: ", err)
	}

	err = env.SetGeometry(defaultGeometry)
	if err != gmdbx.ErrSuccess {
		log.Fatal("set geometry failed: ", err)
	}

	err = env.SetOption(gmdbx.OptTxnDpLimit, 10240)
	if err != gmdbx.ErrSuccess {
		log.Fatal("set txn dp limit failed: ", err)
	}

	err = env.SetOption(gmdbx.OptMaxReaders, 10240)
	if err != gmdbx.ErrSuccess {
		log.Fatal("set max readers failed: ", err)
	}

	err = env.Open(path, defaultFlags, 0664)
	if err != gmdbx.ErrSuccess {
		log.Fatal("open db failed: ", err)
	}

	// init sub dbs
	tx := &gmdbx.Tx{}
	if err = env.Begin(tx, gmdbx.TxReadWrite); err != gmdbx.ErrSuccess {
		log.Fatal("open tx failed: ", err)
	}
	defer tx.Commit()

	db.users, err = tx.OpenDBI(SubDBUsers, gmdbx.DBCreate|gmdbx.DBIntegerKey)
	if err != gmdbx.ErrSuccess {
		log.Fatal("open users db failed: ", err)
	}

	db.permits, err = tx.OpenDBI(SubDBPermits, gmdbx.DBCreate|gmdbx.DBIntegerKey)
	if err != gmdbx.ErrSuccess {
		log.Fatal("open users db failed: ", err)
	}

	db.rels, err = tx.OpenDBI(SubDBRels, gmdbx.DBCreate|gmdbx.DBIntegerKey|gmdbx.DBDupSort)
	if err != gmdbx.ErrSuccess {
		log.Fatal("open rels db failed: ", err)
	}

	db.convsUsers, err = tx.OpenDBI(SubDBConvsUsers, gmdbx.DBCreate|gmdbx.DBIntegerKey|gmdbx.DBDupSort)
	if err != gmdbx.ErrSuccess {
		log.Fatal("open convs db failed: ", err)
	}

	db.convsMesgs, err = tx.OpenDBI(SubDBConvsMesgs, gmdbx.DBCreate|gmdbx.DBIntegerKey|gmdbx.DBDupSort)
	if err != gmdbx.ErrSuccess {
		log.Fatal("open convs db failed: ", err)
	}

	db.mesgs, err = tx.OpenDBI(SubDBMesgs, gmdbx.DBCreate|gmdbx.DBIntegerKey)
	if err != gmdbx.ErrSuccess {
		log.Fatal("open mesgs db failed: ", err)
	}

	db.genv, err = tx.OpenDBI(SubDBEnv, gmdbx.DBCreate)
	if err != gmdbx.ErrSuccess {
		log.Fatal("open freinds db failed: ", err)
	}

	mesgeq, gerr := mq.NewMessageQueue(mqPath, 256<<20, 0)
	if gerr != nil {
		log.Fatal("open message queue failed: ", gerr)
	}
	db.mesgeq = mesgeq
	db.seqm = make(map[string]*atomic.Uint64)

	// init message id
	mesgid := gmdbx.String(&SeqmMesgs)
	mesgval := gmdbx.Val{}
	err = tx.Get(db.genv, &mesgid, &mesgval)
	if err != gmdbx.ErrSuccess && err != gmdbx.ErrNotFound {
		log.Fatal("init message id failed: ", err)
	}

	if err == gmdbx.ErrNotFound {
		db.seqm[SeqmMesgs] = atomic.NewUint64(0)
	} else {
		db.seqm[SeqmMesgs] = atomic.NewUint64(mesgval.U64())
		slog.Info("message id", "message id", mesgval.U64())
	}

	// init conversation id
	convid := gmdbx.String(&SeqmConvs)
	convval := gmdbx.Val{}

	err = tx.Get(db.genv, &convid, &convval)
	if err != gmdbx.ErrSuccess && err != gmdbx.ErrNotFound {
		log.Fatal("init conversation id failed: ", err)
	}
	if err == gmdbx.ErrNotFound {
		db.seqm[SeqmConvs] = atomic.NewUint64(0)
	} else {
		db.seqm[SeqmConvs] = atomic.NewUint64(convval.U64())
		slog.Info("conversation id", "conversation id", convval.U64())
	}

	// remove message which is consumed every 30 seconds
	go func(ctx context.Context) {
		slog.Info("start remove message goroutine")
		tick := time.Tick(30 * time.Second)
		for {
			select {
			case <-ctx.Done():
				db.mesgeq.Close()
				slog.Info("mq closed")
				return
			case <-tick:
				err := db.mesgeq.DeleteConsumedMessages()
				if err != nil {
					slog.Error("delete consumed messages failed: ", "err", err)
				}
			}
		}
	}(ctx)

	// batch wirte data into database
	go func(ctx context.Context) {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		var (
			mi types.MesgInfo
		)

		for {
			select {
			case <-ctx.Done():
				// save sequence
				msq := db.seqm[SeqmMesgs].Load()
				mesgval := gmdbx.U64(&msq)
				wtx, err := db.Wtx()
				if err != nil {
					slog.Error("save sequence tx failed: ", "err", err)
				}

				mkey := gmdbx.String(&SeqmMesgs)
				xerr := wtx.Put(db.genv, &mkey, &mesgval, gmdbx.PutUpsert)
				if xerr != gmdbx.ErrSuccess {
					slog.Error("save message sequence failed: ", "err", xerr)
				}

				slog.Info("message sequence", "msg id", msq)

				conv := db.seqm[SeqmConvs].Load()
				convval := gmdbx.U64(&conv)
				ckey := gmdbx.String(&SeqmConvs)
				xerr = wtx.Put(db.genv, &ckey, &convval, gmdbx.PutUpsert)
				if xerr != gmdbx.ErrSuccess {
					slog.Error("save conv sequence failed: ", "err", xerr)
				}

				slog.Info("convs sequence", "key", SeqmConvs, "convs id", conv)

				werr := wtx.Commit()
				if werr != gmdbx.ErrSuccess {
					slog.Error("commit sequence failed: ", "err", werr)
				}

				slog.Info("sequence saved")
				// close database
				db.Close()
				slog.Info("database closed")
				return
			default:
				mesgs, err := db.mesgeq.PopAll()
				if err != nil {
					slog.Error("pop message failed: ", "err", err)
					continue
				}

				if len(mesgs) == 0 {
					time.Sleep(200 * time.Millisecond)
					continue
				}

				wtx, err := db.Wtx()
				if err != nil {
					slog.Error("begin write tx failed: ", "err", err)
					continue
				}

				// write data into database
				for _, mesg := range mesgs {
					err := sonic.Unmarshal(mesg.Data, &mi)
					if err != nil {
						slog.Error("unmarshal message failed: ", "err", err)
						continue
					}

					handler, ok := mesgHandlers[mi.MsgType]
					if !ok {
						slog.Error("unknown message type: ", "type", mi.MsgType)
						continue
					}

					err = handler(&mi, wtx)
					if err != nil {
						slog.Error("handle message failed: ", "err", err)
						continue
					}
				}
				wtx.Commit()
				slog.Info("message processed", "count", len(mesgs))
			}
		}
	}(ctx)

	return db
}

func (db *DB) Close() error {

	err := db.env.CloseDBI(db.users)
	if err != gmdbx.ErrSuccess {
		slog.Error("close users db failed: ", "err", err)
		return fmt.Errorf("close users db failed: %v", err)
	}

	err = db.env.CloseDBI(db.permits)
	if err != gmdbx.ErrSuccess {
		slog.Error("close permits db failed: ", "err", err)
		return fmt.Errorf("close permits db failed: %v", err)
	}

	err = db.env.CloseDBI(db.rels)
	if err != gmdbx.ErrSuccess {
		slog.Error("close rels db failed: ", "err", err)
		return fmt.Errorf("close rels db failed: %v", err)
	}

	err = db.env.CloseDBI(db.convsUsers)
	if err != gmdbx.ErrSuccess {
		slog.Error("close convs users db failed: ", "err", err)
		return fmt.Errorf("close convs users db failed: %v", err)
	}

	err = db.env.CloseDBI(db.convsMesgs)
	if err != gmdbx.ErrSuccess {
		slog.Error("close convs mesgs db failed: ", "err", err)
		return fmt.Errorf("close convs mesgs db failed: %v", err)
	}

	err = db.env.CloseDBI(db.mesgs)
	if err != gmdbx.ErrSuccess {
		slog.Error("close mesgs db failed: ", "err", err)
		return fmt.Errorf("close mesgs db failed: %v", err)
	}

	err = db.env.CloseDBI(db.genv)
	if err != gmdbx.ErrSuccess {
		slog.Error("close global env db failed: ", "err", err)
		return fmt.Errorf("close global env db failed: %v", err)
	}

	err = db.env.Close(false)
	if err != gmdbx.ErrSuccess {
		slog.Error("close db failed: ", "err", err)
		return fmt.Errorf("close db failed: %v", err)
	}

	return err
}

func (d *DB) Wtx() (*gmdbx.Tx, error) {
	tx := &gmdbx.Tx{}
	err := d.env.Begin(tx, gmdbx.TxReadWrite)
	if err != gmdbx.ErrSuccess {
		return nil, fmt.Errorf("begin tx failed %v", err)
	}
	return tx, nil
}

func (d *DB) Rtx() (*gmdbx.Tx, error) {
	tx := &gmdbx.Tx{}
	err := d.env.Begin(tx, gmdbx.TxReadOnly)
	if err != gmdbx.ErrSuccess {
		return nil, fmt.Errorf("begin tx failed %v", err)
	}
	return tx, nil
}
