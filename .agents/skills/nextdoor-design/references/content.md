# Content design — operational rules

Full reference: `docs/design/nextdoor-content.md`.

## Three registers

| Register | Shape | maulid.id |
| --- | --- | --- |
| Reader | first person, present, greeting/question, explicit ask | short, warm, not performative |
| Staff | institutional, plain | `Dikelola staf maulid.id`, no “I” |
| Business/sponsored | offer + benefit + CTA | never |

Nextdoor separates registers by sentence shape, not only by a label. Keep that property.

## Post mechanics

- Ordinary posts have no headline field; the first sentence is the headline.
- Publisher content has a headline + deck.
- Truncation is inline: “… see more” / `Selengkapnya`.
- One attachment per post.
- Empty comment line: `Belum ada komentar` / `Jadi yang pertama`.

## Dates

| Nextdoor | maulid.id |
| --- | --- |
| `2 hr ago` | `2 jam lalu` |
| `1d` | `1 hari lalu` |
| `Jul 27` | `27 Jul` |
| `11:00 PM` | `23.00 WIB` |
| `2,529 Neighbors` | `2.529 warga` |
| `Edited` | `Diedit` |

Relative time in the feed; absolute date + `WIB` in detail/events. No distance or weather
without geodata.

## Canonical strings

```text
Belum divalidasi
Menunggu peninjauan staf
Sudah divalidasi
Kabar dihapus
Pengaturan tersimpan
Laporan terkirim        (only when reports exist)
Belum ada kabar
Belum ada kabar yang cocok
Kabar tidak dapat dimuat.
Coba lagi
```

Rules: sentence case, no terminal period except prose sentences, one verb per outcome, one
string table. Do not let each feature invent its own toast.

## Validation copy

```text
Belum divalidasi
Akunmu belum divalidasi. Staf meninjau akun sebelum penerbitan dibuka.
[Lihat status akun] [Cara validasi]

Menunggu peninjauan staf
Pengajuanmu sedang ditinjau. Kami akan memperbarui status akun ini.

Sudah divalidasi
Akunmu sudah divalidasi. Kamu dapat menerbitkan kabar.
[Tulis kabar]
```

## Voice

- Prefer `warga`, `jamaah` or `komunitas` over `tetangga`.
- `Assalamualaikum` only if the editorial voice already uses it; otherwise `Halo, warga`.
- No emoji density, no exclamation in every sentence, no copied US idioms.
- No fake urgency, engagement bait, fake social proof or unsupported “verified” claims.
- All labels, dates, empty states and errors in Indonesian.

## Rules page

Four short rules, linked from the composer and contribution section:

1. Kabar harus benar dan dapat dipertanggungjawabkan.
2. Hormati warga dan jangan menyerang pribadi.
3. Jangan menyebar ujaran kebencian atau hoaks.
4. Kabar terkait Maulid dan komunitas; bukan iklan.
