# Public landing-page direction

## Subject inventory and argument

The real material is the name maulid.id; Indonesian Maulid Rasulullah community news; a public reading surface; accounts that can interact before they are validated; validated accounts that publish; staff review of accounts; news posts; shared public/mobile API data; Indonesian interface copy; and the constraint that only approved community contributors can publish. No client photographs, confirmed event schedule, business statistics, testimonials, or established visual identity were supplied.

Audience: Indonesian readers looking for Maulid news and community contributors wanting to share it. The page should show what the community shares and explain who may publish. Its evidence is the actual public post feed and the actual account validation flow.

Page argument:

- Maulid news has a common reading place: the masthead, short explanation and real feed link.
- Readers can find and respond to posts: the API-backed feed, a working text filter and account interaction.
- Publishing requires account validation: a short explanation that matches the backend and account status.
- The account and staff surfaces are reachable: actual navigation and protected backoffice links, with honest image credits in the close.

## Dealt direction and decisions

Used `design-taste-frontend` from `/home/dani/.agents/skills/design-taste-frontend`.

- Dealt three options: `page-is-the-object`, `editorial-report`, `contact-sheet`.
- Chose `editorial-report` (The printed annual), seed **979061848**, because a community news feed benefits from a clear editorial reading axis and masthead; neither fictional invitations nor a large archive of project photographs exists.
- The complete machine-readable direction is in `.taste-skill/direction.json`.
- Retain the editorial axis, photographic masthead, coherent type and flat palette, phone-specific composition, short headlines and masked-headline signature motion.
- Redrew the `list` section to `list-two-ended-rows`: the original day-tab agenda assumes a dated program that was not supplied. Posts form a real list with titles and metadata. Let bodies expand naturally for readable content and error/empty states.
- Where the sheet conflicts with the house rules, follow the rules: no accent-styled headline word, no eyebrow/index furniture, no pointless near-empty full-viewport footer, no decorative tiles without content or honest imagery, and readable body text.

Idea sentence: **On this page, maulid.id's distinction between reading together and publishing after validation becomes an edited community bulletin, with the actual posts occupying its reading column.** The approach belongs to this product's account/validation constraint, but the editorial metaphor is shared with other community-news products; do not overstate its originality in the final critique.

Live thing: the API-backed public post feed can be filtered by a visitor's words and interacted with after signing in. It uses actual post content, works with a thumb, and leaves the unfiltered readable list available when the filter is empty. Empty installations show an honest empty state; any seeded demo posts must be explicitly labeled and opt-in.

## Copy deck

- Brand: `maulid.id`
- Navigation: `Kabar`, `Kontribusi`, `Akun`, `Backoffice` where applicable.
- Hero: `Kabar Maulid kita` (17 characters, safe within two lines).
- Supporting sentence: `Baca kabar Maulid Rasulullah dari komunitas Indonesia, lalu ikut berbagi lewat akun yang sudah divalidasi.`
- Primary action: `Baca kabar`, a real filled control visible on the first screen that reaches the feed.
- Feed heading: `Kabar komunitas`.
- Search label: `Cari kabar`; empty result: `Belum ada kabar yang cocok. Coba kata lain.`
- Empty installation: `Belum ada kabar. Akun yang sudah divalidasi dapat membagikan kabar pertama.`
- Contribution heading: `Ikut berbagi`.
- Contribution explanation: `Setiap akun dapat berinteraksi dengan kabar. Untuk menerbitkan kabar, akunmu perlu divalidasi oleh staf.`
- Canonical user-facing states: `Belum divalidasi`, `Sudah divalidasi`, with their documented API values used consistently.
- Footer: brand, actual navigation, brief community purpose and actual image credits; no invented contact address or social profiles.

## Imagery and verification

Use four honestly contextual, licensed photographs from the Indonesian mosque/community environment if available, downloaded locally with their creators, source pages and licenses documented and credited. Do not represent them as confirmed maulid.id events or platform members. The dealt generic hiking/portrait photographs do not establish this subject and should not be used. Ask the research worker for sources; if no usable source is available, follow the skill's documented fallback honestly and report the limitation.

The implementation worker must run the skill's lint and real-browser inspector against the public page, inspect desktop/tablet/mobile screenshots, write the specified critique, fix concrete problems, and verify reduced-motion behavior. At most three critique rounds. Record actual scores and limitations; no invented measurements. Run ordinary responsive/auth/form checks for backoffice separately, since the public-page art direction is not a dashboard style mandate.
