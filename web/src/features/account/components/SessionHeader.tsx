import { useLocation } from "react-router-dom";
import { SiteHeader } from "@/shared/components/SiteHeader";
import { useSession } from "../session";

export function SessionHeader() {
  const { account, status } = useSession();
  const location = useLocation();
  return (
    <SiteHeader
      account={account}
      loading={status === "loading"}
      overlay={location.pathname === "/"}
    />
  );
}
