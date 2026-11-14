export function Masthead() {
  return (
    <section className="masthead" aria-labelledby="hero-title">
      <div className="masthead__band" aria-hidden="true">
        <figure className="masthead__cell masthead__cell--a">
          <img
            src="/images/masjid-istiqlal-jakarta.webp"
            width="1280"
            height="856"
            alt=""
            fetchPriority="high"
            decoding="async"
          />
        </figure>
        <figure className="masthead__cell masthead__cell--b">
          <img
            src="/images/masjid-baiturrahman-aceh.webp"
            width="1600"
            height="1202"
            alt=""
            fetchPriority="high"
            decoding="async"
          />
        </figure>
        <figure className="masthead__cell masthead__cell--c">
          <img
            src="/images/masjid-menara-kudus.webp"
            width="1600"
            height="1066"
            alt=""
            loading="lazy"
            decoding="async"
          />
        </figure>
        <figure className="masthead__cell masthead__cell--d">
          <img
            src="/images/peringatan-maulid-nabi-2022-di-masjid-nurul.webp"
            width="1600"
            height="1200"
            alt=""
            loading="lazy"
            decoding="async"
          />
        </figure>
        <div className="masthead__tile masthead__tile--blue" />
        <div className="masthead__tile masthead__tile--red" />
        <div className="masthead__tile masthead__tile--black" />
      </div>

      <div className="masthead__copy">
        <p className="masthead__standfirst">
          Baca kabar Maulid Rasulullah dari komunitas Indonesia, lalu ikut berbagi lewat akun yang
          sudah divalidasi.
        </p>
        <h1 className="masthead__title" id="hero-title">
          <span className="masthead__line">
            <span className="masthead__line-inner">Kabar</span>
          </span>
          <span className="masthead__line">
            <span className="masthead__line-inner">Maulid kita</span>
          </span>
        </h1>
        <a className="button button--primary masthead__action" href="/#kabar">
          Baca kabar
        </a>
      </div>

      <p className="masthead__caption">
        Foto suasana: Masjid Istiqlal, Masjid Baiturrahman, Menara Kudus, dan peringatan Maulid di
        Jombang. Bukan dokumentasi kegiatan maulid.id.
      </p>
    </section>
  );
}

import "./Masthead.css";
