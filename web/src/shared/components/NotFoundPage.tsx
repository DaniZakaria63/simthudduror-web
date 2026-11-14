import { Link } from "react-router-dom";

export function NotFoundPage() {
  return (
    <section className="not-found" aria-labelledby="not-found-title">
      <h1 id="not-found-title">Halaman tidak ditemukan</h1>
      <p>Alamat ini tidak tersedia. Kembali ke kabar komunitas atau halaman akun.</p>
      <p>
        <Link className="button button--primary" to="/">
          Buka beranda
        </Link>
      </p>
    </section>
  );
}

import "./NotFoundPage.css";
