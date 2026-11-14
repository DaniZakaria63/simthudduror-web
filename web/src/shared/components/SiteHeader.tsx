import { Link } from "react-router-dom";

type HeaderAccount = {
  name: string;
  role: string;
};

type SiteHeaderProps = {
  account: HeaderAccount | null;
  loading?: boolean;
  overlay?: boolean;
};

export function SiteHeader({ account, loading = false, overlay = false }: SiteHeaderProps) {
  const action =
    loading ? null : account?.role === "staff" ? (
      <Link className="site-header__action" to="/backoffice">
        Backoffice
      </Link>
    ) : account ? (
      <Link className="site-header__action" to="/akun">
        {account.name}
      </Link>
    ) : (
      <Link className="site-header__action" to="/masuk">
        Masuk
      </Link>
    );

  return (
    <header className={`site-header ${overlay ? "site-header--overlay" : ""}`.trim()}>
      <Link className="site-header__brand" to="/" aria-label="maulid.id, beranda">
        maulid.id
      </Link>
      <nav className="site-header__nav" aria-label="Navigasi utama">
        <a href="/#kabar">Kabar</a>
        <a href="/#kontribusi">Kontribusi</a>
        <Link to="/akun">Akun</Link>
      </nav>
      <div className="site-header__end">{action}</div>
    </header>
  );
}

import "./SiteHeader.css";
