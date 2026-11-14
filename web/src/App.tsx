import { useLocation, Route, Routes } from "react-router-dom";
import { AccountPage, AccountProvider, LoginPage, RegisterPage, SessionHeader } from "@/features/account";
import { BackofficePage } from "@/features/backoffice";
import { LandingPage } from "@/features/landing";
import { NotFoundPage } from "@/shared/components/NotFoundPage";
import { SiteFooter } from "@/shared/components/SiteFooter";
import { useReveal } from "@/shared/hooks/useReveal";

function AppChrome() {
  const location = useLocation();
  useReveal(location.pathname);

  return (
    <>
      <a className="skip-link" href="#main">
        Lewati ke konten
      </a>
      <SessionHeader />
      <main id="main" className="app-main">
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/masuk" element={<LoginPage />} />
          <Route path="/daftar" element={<RegisterPage />} />
          <Route path="/akun" element={<AccountPage />} />
          <Route path="/backoffice" element={<BackofficePage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </main>
      <SiteFooter />
    </>
  );
}

export default function App() {
  return (
    <AccountProvider>
      <AppChrome />
    </AccountProvider>
  );
}

import "./App.css";
