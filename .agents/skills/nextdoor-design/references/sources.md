# Sources and capture notes

## Method

- Product: Nextdoor.
- Library: Mobbin, via the Mobbin MCP server (`https://api.mobbin.com/mcp`).
- Platforms: web and iOS.
- Queries covered the home feed, post detail, composer, groups, events, search,
  notifications, profile, onboarding, moderation and seasonal modules; no Nextdoor seasonal
  capture was found.
- Every capture was inspected as an image. Product names and URLs are preserved below.
- Task intent sent with every Mobbin call: “Extract Nextdoor design system and screen
  patterns for maulid.id”.
- Capture date: 2026-11-15.

The complete captured screen list is in `docs/design/nextdoor-screen-inventory.md`. This
file keeps the source rules, the key URLs and the measurement notes.

## What the evidence supports

- Screen composition, component presence, labels, hierarchy and visible states.
- Approximate proportions, spacing and colour values sampled from the images.
- Interaction order where a flow capture shows it.

## What it does not support

- Exact production design tokens, CSS, source code or font files.
- Claims that a pattern improved retention, trust or conversion.
- Permission to reuse Nextdoor's logo, brand assets, illustrations, photography, copy,
  advertising formats or exact compositions.
- Complete coverage of every Nextdoor screen. Mobbin is a capture library, not the product.

## Confirmed Nextdoor captures

### Web

| Capture | URL | What it shows |
| --- | --- | --- |
| Web home (landscaping search) | <https://mobbin.com/screens/03d5f316-fcd3-4389-9cc0-a5406574beca> | three-zone shell, left nav, feed, right rail, Post CTA, ads |
| Web home (house-cleaner search) | <https://mobbin.com/screens/50649f88-3142-465e-ab6d-037207fd3f78> | same shell, feed card, invite card, business CTA |
| Web home — mixed feed | <https://mobbin.com/screens/ed19ca8f-d48e-4706-8f8c-365b5707e096> | listing + publisher news, rail invite/business cards |
| Web home — long-form post | <https://mobbin.com/screens/d6cea51b-2236-4bc1-8d8f-1d44a77d320d> | long body, reactions, inline notice |
| Web home — onboarding checklist | <https://mobbin.com/screens/e42a20e9-11ab-43ff-bf03-264aec66f035> | “Get started” card + progress |
| Community rail | <https://mobbin.com/screens/efe11897-4db8-46ec-b5f1-e6d8fa164e53> | neighbourhood stats |
| Local publisher sheet | <https://mobbin.com/screens/d748f96c-1850-44c9-b0c3-2f38894f97b5> | publisher identity, article, comments |
| Profile | <https://mobbin.com/screens/93fa03f3-e6dd-49da-94ff-64942a7af20e> | banner, avatar, Message pill, posts, skeleton |
| Events near you | <https://mobbin.com/screens/3a2e766c-4903-421d-a04f-5d72155d6dbe> | 3-up event grid, RSVP pills |
| For Sale grid | <https://mobbin.com/screens/858c33ce-f874-4c54-9619-99ed7788e654> | 6-up marketplace grid |
| Groups directory | <https://mobbin.com/screens/2b33c741-7393-4cd4-91b3-15bd9c7b0775> | list + Your groups rail |
| Group page | <https://mobbin.com/screens/80a03c7a-f8fc-4f8e-9a78-4c01eeec5899> | cover, join state, welcome carousel, pinned post |
| Faves answer | <https://mobbin.com/screens/f0a1d80a-8b90-4859-9304-546d24871cac> | AI answer, sources, feedback |
| Sign up | <https://mobbin.com/screens/e2813a0f-740c-4c55-8e02-a2d139368d97> | logged-out shell, OAuth, email form |
| Business dashboard | <https://mobbin.com/screens/e3866f6f-4a88-4e2c-bacd-112b06ce2e55> | business nav, setup checklist |

### iOS

| Capture | URL | What it shows |
| --- | --- | --- |
| Home feed | <https://mobbin.com/screens/251e08e3-4688-4ea7-be88-557ec9871142> | location header, post + event card, local-news rail, tab bar, Post FAB |
| Local Deals explanation | <https://mobbin.com/screens/feb2cdda-699f-4be5-a11a-244612039997> | local-deals explanation page; **not** a seasonal/Ramadan card |
| Search results | <https://mobbin.com/screens/6d9b9a65-a2bd-4f67-b0e4-4e7ba3f0068a> | tabs, filter dropdowns, text-first results |
| Get started onboarding | <https://mobbin.com/screens/4dd17697-a94f-4a82-a80e-49f770a952dc> | growth cards, progress, sponsored post |
| For Sale & Free + local news | <https://mobbin.com/screens/54cefa73-8bc8-4c81-add3-ed34b1f570fe> | listing attachment, local-news rail |
| Text + image + sponsored | <https://mobbin.com/screens/7f3764d0-da67-46f0-b6c6-e49dd4fb1484> | photo post, sponsored card, bottom bar |
| Mixed feed + sponsored | <https://mobbin.com/screens/ed9fc5ab-889d-45e3-81bf-015ad6b46c69> | text recommendations, sponsored card |
| Local news + Target ad | <https://mobbin.com/screens/2aeb4bc9-a474-4696-92b6-5f1b3994a0a9> | news card + neighbour reaction, sponsored card |
| For Sale & Free list | <https://mobbin.com/screens/e1316849-84ec-4fa3-914c-c8be49a854c3> | listing-heavy feed |
| Dark mode | <https://mobbin.com/screens/d178d1a1-d1fe-4cc4-b19e-24dc5f6fd108> | dark canvas/cards, sponsored card |
| Poll + event | <https://mobbin.com/screens/1aae4c0b-e83c-4e67-9c2c-e5a11cff8698> | poll options, event date block |
| Local news + neighbour reaction | <https://mobbin.com/screens/ed239926-f476-46ca-87ec-e897ff20028b> | publisher card, neighbour question, reactions |
| Group feed | <https://mobbin.com/screens/9fa1c528-5530-4f11-bc24-0244ea7279fe> | groups module, posts, poll/place attachments |
| News + neighbour commentary flow | <https://mobbin.com/flows/f3c0a105-1443-4497-9f1c-20aa376762aa?tab=screens> | publisher lockup, headline, photo, neighbour reaction |
| Composer flow | <https://mobbin.com/flows/e258bc98-834d-4cce-b62a-cbc96c477653?tab=screens> | author identity, audience, intent cards |
| Post overflow flow | <https://mobbin.com/flows/b278d4f1-98e7-4915-91c2-264f86cff9c6?tab=screens> | hide, mute, report, share, bookmark with consequences |
| Community guidelines | <https://mobbin.com/screens/7546ad42-8bdd-4c1c-894c-5ec39388788c> | warm opening, filters, numbered rules |

Observed bottom navigation labels: `Home`, `Search`, `For Sale`, `Faves`, and a raised green
`Post` FAB. The active destination is green; the FAB is a circular compose action at the
right, not a centre tab.

Additional captures are cited inline in `web-app.md`, `mobile-app.md`, `components.md` and
`patterns.md`. If a capture is not listed there, it was not used as evidence.

## Note: no Nextdoor marketing-site sections

Mobbin returned no Nextdoor website sections. The captures are the web app and iOS app. Do
not describe a Nextdoor marketing site, colour palette or campaign layout as evidence; the
public site was not part of this corpus.

## Correction: no Nextdoor seasonal module

An earlier draft of the broader Mobbin research misread
`https://mobbin.com/screens/feb2cdda-699f-4be5-a11a-244612039997` as a “Ramadan Mubarak”
card. The Mobbin page title is **Nextdoor iOS Local Deals — local deals explanation page**.
It is not seasonal evidence. Searches for `Nextdoor Ramadan`, `Nextdoor seasonal holiday
card` and related queries returned no Nextdoor seasonal module.

Do not cite a Nextdoor Ramadan card. The Maulid seasonal module in the maulid.id design is
an editorial decision, not a Nextdoor-derived pattern.

## Visual measurements

- Web capture dimensions: 1920×1320.
- Sampled canvas: `#F8F9FA`.
- Sampled surface: `#FFFFFF`.
- Sampled ink: `#2F2F2F`.
- Sampled border: `#EBEDED`.
- Sampled primary green: `#188050`; darker `#107040`.
- Sampled dark CTA: `#222E46`.
- Web proportions (measured): left nav ~17%, feed ~38–39%, right rail ~19%; remaining width
  is margin. Feed:rail ≈2:1.
- Values not sampled are labelled **estimated** in the design documents.

## Copyright and use

The captures remain Nextdoor's and Mobbin's property. This repository records observations
and URLs only; it does not vendor the images, logo, illustrations, copy or other assets.
Future UI work must adapt patterns, not reproduce Nextdoor's protected design.
