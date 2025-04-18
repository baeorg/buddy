import type { Route } from "./+types/login";
import { Login } from "../page/login";
export function meta({}: Route.MetaArgs) {
  return [{ title: "Login" }, { name: "description", content: "Login" }];
}

export default function LoginRoute() {
  return <Login />;
}
