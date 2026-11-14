import { Link } from "react-router-dom";

export function SiteFooter() {
  return (
    <footer className="site-footer" id="tentang">
      <div className="site-footer__columns">
        <div data-reveal>
          <p className="site-footer__strong">maulid.id</p>
          <p className="site-footer__note">
            Ruang baca kabar Maulid Rasulullah dari komunitas Indonesia. Setiap akun dapat
            berinteraksi; penerbitan dilakukan akun yang sudah divalidasi staf.
          </p>
        </div>
        <nav aria-label="Navigasi footer" data-reveal data-reveal-delay="1">
          <p className="site-footer__strong">Jelajahi</p>
          <ul>
            <li>
              <a href="/#kabar">Kabar</a>
            </li>
            <li>
              <a href="/#kontribusi">Kontribusi</a>
            </li>
            <li>
              <Link to="/akun">Akun</Link>
            </li>
            <li>
              <Link to="/backoffice">Backoffice staf</Link>
            </li>
          </ul>
        </nav>
        <div data-reveal data-reveal-delay="2">
          <p className="site-footer__strong">Foto</p>
          <p className="site-footer__note">
            Foto suasana masjid dan peringatan Maulid di Indonesia, bukan dokumentasi kegiatan
            maulid.id.
          </p>
        </div>
        <div className="site-footer__credits" data-reveal data-reveal-delay="3">
          <p className="site-footer__strong">Kredit gambar</p>
          <ul>
            <li>
              <a href="https://commons.wikimedia.org/wiki/File:Masjid_Istiqlal_Jakarta.jpg">
                Masjid Istiqlal Jakarta
              </a>{" "}
              oleh AR El Syarif,{" "}
              <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>, Wikimedia
              Commons.
            </li>
            <li>
              <a href="https://commons.wikimedia.org/wiki/File:Masjid_Baiturrahman_Aceh.jpg">
                Masjid Baiturrahman Aceh
              </a>{" "}
              oleh ROFIQOHHASAN,{" "}
              <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>, Wikimedia
              Commons.
            </li>
            <li>
              <a href="https://commons.wikimedia.org/wiki/File:Masjid_Menara_Kudus.jpg">
                Masjid Menara Kudus
              </a>{" "}
              oleh PL09Puryono,{" "}
              <a href="https://creativecommons.org/publicdomain/zero/1.0/deed.en">CC0</a>, Wikimedia
              Commons.
            </li>
            <li>
              <a href="https://commons.wikimedia.org/wiki/File:Peringatan_Maulid_Nabi_2022_di_Masjid_Nurul_Huda_Pulo_Lor-13.jpg">
                Peringatan Maulid Nabi 2022 di Masjid Nurul Huda Pulo Lor
              </a>{" "}
              oleh Indonesiagood,{" "}
              <a href="https://creativecommons.org/licenses/by/4.0/">CC BY 4.0</a>, Wikimedia
              Commons.
            </li>
            <li>
              Huruf EB Garamond oleh The EB Garamond Project Authors,{" "}
              <a href="https://openfontlicense.org/">SIL Open Font License 1.1</a>.
            </li>
          </ul>
        </div>
      </div>
      <p className="site-footer__legal">
        maulid.id · Kabar komunitas, bukan jadwal resmi penyelenggara.
      </p>
      <h2 className="site-footer__wordmark" aria-label="maulid.id">
        maulid.id
      </h2>
    </footer>
  );
}

import "./SiteFooter.css";
