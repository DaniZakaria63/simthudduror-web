# Nextdoor content design and tone

Evidence: the captures in `nextdoor-screen-inventory.md` and the microcopy table in the
cross-cutting research notes. Short UI strings are quoted as evidence; longer body copy is
paraphrased.

## 1. Three voice registers

Nextdoor separates voices by sentence shape, not only by a label. A reader who ignores the
badge can still tell who is speaking.

| Register | Shape | Example evidence |
| --- | --- | --- |
| Neighbour | first person, present tense, greeting or question, explicit ask | “Hi neighbours! I'm Alex Smith — new to the neighbourhood…” |
| Publisher | third person, headline + deck + verb | “West Menlo Park store opens despite penalties…” + “Read article” |
| Business | offer + benefit + CTA, second person imperative | “Expert HVAC Services” + “Schedule Service Now” |

**maulid.id:**

- Reader/account: first person, short, warm but not performative.
- Staff: institutional, plain, no first-person “I”.
- No business register. There is no sponsored voice.
- The validation label is part of the identity line, not a separate marketing claim.

## 2. Headline and body mechanics

- Ordinary posts have **no headline field**. The first sentence is the headline; it is the
  same size as the body on web, slightly heavier on iOS.
- Publisher posts have a real headline and deck.
- Truncation is inline: a paragraph ends with “… see more” as a text run, not a separate
  button.
- One attachment per post: image, event, listing, poll or video — never two.
- Lists carry their own empty line (`No comments` / `Be the first to comment`).
- Poll results show real percentages and counts.

**maulid.id:** body-only today. When a title field is added, keep the body excerpt and a
deterministic “Selengkapnya” control; do not truncate the headline before the dek.

## 3. Microcopy inventory

| Job | Nextdoor observed | maulid.id rule |
| --- | --- | --- |
| Success | `Your settings have been saved` / `Successfully updated your profile.` | one string per outcome: `Pengaturan tersimpan` |
| Delete | `Post deleted` / `Discussion closed.` | `Kabar dihapus` (no period) |
| Report | `Report submitted.` | `Laporan terkirim` (only when reports exist) |
| Mute | `You Muted MassLive.` | `Kabar dari … disembunyikan` |
| Invite | `Invite sent to 1 person.` | `Undangan terkirim` (future) |
| Empty | `No updates` / `All quiet in the neighborhood.` | `Belum ada kabar` / `Belum ada kabar yang cocok` |
| No comments | `No comments` / `Be the first to comment` | `Belum ada komentar` / `Jadi yang pertama` |
| Error | field-level red text | `Kabar tidak dapat dimuat.` + `Coba lagi` |

**Tone drift warning:** Nextdoor mixes `saved`/`updated`, `Post deleted` (no period) and
`Discussion closed.` (period). This is what happens when each feature writes its own toast.
maulid.id must keep one string table: sentence case, no terminal period, one verb per
outcome.

## 4. Dates, times and numbers

| Nextdoor | maulid.id |
| --- | --- |
| `25m`, `2 hr ago`, `4h` | `25 mnt lalu`, `2 jam lalu` |
| `1d`, `6 days ago` | `1 hari lalu`, `6 hari lalu` |
| `Jul 27`, `Jul 31, 2025` | `27 Jul`, `31 Juli 2025` |
| `11:00 PM` | `23.00 WIB` |
| `Jul 24, 12:00 AM – 1:00 AM WIB` | `24 Jul, 00.00–01.00 WIB` |
| `2,529 Neighbors` | `2.529 warga` |
| `0.69 mi away` | omit; no geodata |
| `58°F Partly Cloudy` | omit |
| `Edited` | `Diedit` |
| `Local publisher · 16 hr ago` | `Akun tervalidasi staf · 16 jam lalu` |

Rules:

- Relative time in the feed; absolute date + timezone in detail/events.
- Always suffix `WIB` on event times. Indonesia has WIB/WITA/WIT; do not assume one zone
  silently.
- Use Indonesian month names in long dates, short month abbreviations in cards.
- No distance or weather unless the API provides it.

## 5. Empty, loading, error and permission states

| Situation | Nextdoor | maulid.id |
| --- | --- | --- |
| Empty feed, signed out | — | `Belum ada kabar.` + `Kabar pertama akan tampil di sini.` No CTA. |
| Empty feed, validated | `No updates` | `Belum ada kabar.` + `Tulis kabar pertama.` |
| Empty feed, unvalidated | — | `Belum ada kabar.` + `Akunmu belum divalidasi.` + link to Akun |
| Filter/search empty | `No suggestions found` | `Belum ada kabar yang cocok.` + `Hapus filter` |
| Error | — | `Kabar tidak dapat dimuat.` + `Coba lagi` |
| Permission | location pre-prompt | not applicable; do not invent notification priming |
| Success | `Post deleted` | `Kabar dihapus` |
| Validation pending | `Thank you for submitting verification request.` | `Pengajuan validasi diterima.` + review expectation |

Nextdoor's useful rules:

- Illustration only when the user can act (events, chats, alerts); dense operational lists
  get text only.
- An empty state may have **no CTA**; `No updates` is deliberately actionless, which is
  calmer than a permanent “post something!” nag.
- A permission prompt explains the benefit and offers a manual alternative.

## 6. Indonesian voice

- Prefer `warga`, `jamaah` or `komunitas` over `tetangga`; `tetangga` implies a shared
  address that maulid.id does not have.
- Use `Assalamualaikum` only if the editorial voice already uses it; do not paste it onto
  every screen. `Halo, warga` is a neutral alternative.
- Keep greetings warm but not Californian-suburban. No emoji density, no exclamation marks
  in every sentence.
- Do not translate Nextdoor idioms literally (`Hi neighbour!`, `All quiet in the
  neighbourhood`).
- `Belum divalidasi` / `Sudah divalidasi` are canonical; do not vary with synonyms.
- Do not imply official event coverage, religious authority or government endorsement.

## 7. Validation and account copy

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

- The account screen explains what validation enables: publishing kabar.
- The composer shows the state at the top; the publish button is disabled with the reason.
- After validation, state the capability, not a generic congratulations.

## 8. Rules copy

Nextdoor's onboarding guidelines use four icon rows:

- Be helpful.
- Lead with respect.
- Do no harm.
- All are welcome.

**maulid.id** should publish four short rules, for example:

1. Kabar harus benar dan dapat dipertanggungjawabkan.
2. Hormati warga dan jangan menyerang pribadi.
3. Jangan menyebar ujaran kebencian atau hoaks.
4. Kabar terkait Maulid dan komunitas; bukan iklan.

Link them from the composer footer and the account/contribution section. Do not gate every
visit with a modal.

## 9. What not to write

- No copied Nextdoor copy, neighbourhood idioms or ad testimonial voice.
- No “verified” claims that the staff review cannot support.
- No urgency (“segera”, “terbatas”) or countdowns.
- No engagement bait (“like kalau setuju”).
- No fake social proof (“ribuan warga membaca”).
- No English-only labels for Indonesian users.
