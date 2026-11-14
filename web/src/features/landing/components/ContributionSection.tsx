import { Link } from "react-router-dom";

export function ContributionSection() {
  return (
    <section className="contribution" id="kontribusi" aria-labelledby="kontribusi-title">
      <div className="contribution__steps">
        <div className="contribution__step contribution__step--1" data-reveal>
          <h2 id="kontribusi-title">Pertama, buat akun.</h2>
          <p>Daftar dengan email, lalu masuk dari web atau aplikasi mobile.</p>
        </div>
        <div className="contribution__step contribution__step--2" data-reveal>
          <h3>Lalu, baca kabar.</h3>
          <p>Setiap akun dapat membaca dan menyukai kabar tanpa menunggu validasi.</p>
        </div>
        <div className="contribution__step contribution__step--3" data-reveal>
          <h3>Minta validasi staf.</h3>
          <p>Staf meninjau akun dari backoffice sebelum memberi akses penerbitan.</p>
        </div>
        <div className="contribution__step contribution__step--4" data-reveal>
          <h3>Setelah itu, terbitkan.</h3>
          <p>Akun yang sudah divalidasi dapat membagikan kabar untuk dibaca bersama.</p>
        </div>
      </div>
      <div className="contribution__action" data-reveal>
        <Link className="contribution__button" to="/daftar">
          Buat akun
        </Link>
        <p>Sudah punya akun? <Link to="/masuk">Masuk di sini</Link>.</p>
      </div>
    </section>
  );
}

import "./ContributionSection.css";
