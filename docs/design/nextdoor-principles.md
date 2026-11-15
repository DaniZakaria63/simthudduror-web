# Nextdoor design principles — evidence and maulid.id application

These principles are distilled from the Nextdoor captures cited in
`../.agents/skills/nextdoor-design/references/sources.md` and the screen notes in
`nextdoor-reference.md`. Each principle states the visible evidence first, then the rule
maulid.id adopts.

## 1. Local context before content

**Evidence.** The web and iOS headers put the neighbourhood name in the first visual
position (for example “Lower Allston”, “Downtown Menlo Park”). The feed is framed as a
place, not an algorithmic stream.

**Rule.** A reader must know what community and scope they are reading before the first
post. maulid.id uses “Kabar komunitas” and can add a region label later; it must not open
with an unexplained global feed.

## 2. The feed is the product

**Evidence.** Home is a single-column stack of posts, with a small set of filters and a
persistent composer. Groups, events, notifications and profile are secondary destinations.
The right rail supports the feed; it does not compete with it.

**Rule.** Spend design effort on the reading item and the transition into publishing.
Do not build a dashboard of modules around the feed.

## 3. Identity before content

**Evidence.** Every post starts with an avatar/initial, a name, a neighbourhood and a
relative time. Publisher items add a publisher lockup and a verification label. Comments
repeat the same identity grammar.

**Rule.** Every maulid.id post shows `author · validation status · time` before the body.
A post without provenance is not a community post.

## 4. Trust is plain language, not decoration

**Evidence.** Nextdoor uses role words (“Local publisher”, “Admin · 2w”, “Sponsored”) and
an address-gated signup whose verification is not shown as a badge on ordinary neighbours.
It does not provide a paid-check mark for everyone; the strongest visible trust signal is a
publisher/agency role label.

**Rule.** Use `Belum divalidasi` / `Sudah divalidasi` in text, with a short explanation on
the account and composer. No paid-check iconography, no trophy, no badge collection.

## 5. One primary action per surface

**Evidence.** Home has the green **Post** action; the composer has **Post**; an event card
has **Interested** or **Join**; the invite card has **Send invitation**. Secondary actions
sit in overflow or quiet buttons.

**Rule.** Each screen has one filled primary action. On the public feed it is “Tulis
kabar” for a validated account, or “Minta validasi” for an unvalidated account. Everything
else is secondary.

## 6. Stack, do not grid

**Evidence.** The home feed is a vertical stack of wide cards. Groups, events and search
results can use compact rows, but the main reading surface is never an equal-weight masonry
grid.

**Rule.** The public feed is one column. Images, when added, are attached to the post, not
used to build a gallery. The editorial masthead may use a composed image band; the feed
does not.

## 7. Progressive disclosure

**Evidence.** Overflow menus hold hide/mute/report/share/bookmark; the composer holds
audience, media and rules; settings hold privacy and notification detail. The first screen
shows only the decision the user needs now.

**Rule.** Put destructive, rare or administrative actions in an overflow menu or a sheet.
Do not put every option on the card. Every overflow item states its consequence.

## 8. The right rail is utility, not content

**Evidence.** The web right rail carries neighbourhood alerts, invitation status, a
business CTA, ads and a chat dock. It is narrower and visually quieter than the feed.

**Rule.** maulid.id's right rail may carry “Tentang komunitas”, account validation status,
staff guidance, or a seasonal Maulid card. It must not become a second feed. On tablet and
phone it collapses below the feed or disappears.

## 9. Honest states beat decorative states

**Evidence.** Empty and error states in the reference set are short, plain and often offer
one recovery action. Filler cards and invented recommendations are rare.

**Rule.** Empty, loading, error and permission states use Indonesian copy and one next
action. Never seed fake posts, fake counts or fake testimonials to fill a screen.

## 10. Calm utility over engagement mechanics

**Evidence.** Reactions are simple glyphs or short labels. Nextdoor does not expose a
global point score, streak or leaderboard on the feed.

**Rule.** Show only real interactions: like, comment (when available), share (when
available). Do not add points, streaks, badges, ranking or “popularity” scores. Do not
optimise for time-on-site.

## 11. The same grammar across web and mobile

**Evidence.** The post card, author row, event card, publisher label and overflow menu keep
their order and meaning on web and iOS; the layout changes, the grammar does not.

**Rule.** A kabar card has the same information order on the public page and the future
mobile app. The API supplies one canonical shape; the UI does not invent platform-specific
labels.

## 12. Accessibility is structural, not a final pass

**Evidence.** Navigation pairs icons with labels, touch targets are generous, and primary
actions are large. The reference still needs independent accessibility verification; a
capture cannot prove contrast or focus behaviour.

**Rule.** Build focus-visible, labels, landmarks, 44px targets and reduced-motion support
into the component, not a later ticket. Verify with keyboard, zoom and a real browser.

## 13. Monetization is separate or absent

**Evidence.** Nextdoor labels sponsored content but gives it the same card geometry as
editorial content; local deals and business CTAs occupy the rail and feed. This is a
warning about ambiguity.

**Rule.** maulid.id has no ads, sponsored cards, paywall or promo countdown at launch. If a
future partnership exists, it must be visually and semantically distinct from kabar, and
it requires an explicit product decision.

## 14. Two primaries, two jobs

**Evidence.** Nextdoor uses green for the social commit (`Post`, `Continue`, FAB) and navy
for the operational commit (`Join`, `Interested`, `Save`, `Next`). Destructive actions use
red. The distinction is consistent enough to read as a rule.

**Rule.** maulid.id uses one green primary for publishing/positive actions and a dark
`--ui-strong` for operational confirmation. Do not use both colours for the same job, and
do not use red for emphasis.

## 15. Commitment stays in the thumb zone

**Evidence.** Onboarding, forms, the composer and report screens pin the primary action to
the bottom, above the safe area, and keep it visible when the keyboard opens.

**Rule.** The publish/save/continue action is reachable with one thumb on mobile and never
hidden behind a scroll. The header can hold a secondary commit only when the form is short.

## 16. Media is optional

**Evidence.** Text-only recommendation posts sit beside photo, poll, event and listing posts
with the same identity and action grammar. A post without an image is still complete.

**Rule.** Do not require an image, and do not render an empty media frame. The text-first
card is the default; media is an enhancement.

## 17. System-native restraint

**Evidence.** Nextdoor uses standard iOS sheets, safe areas, keyboards and back gestures;
custom chrome is limited to the FAB, capsule tab bar and line-art illustrations. The content
carries the colour.

**Rule.** Prefer the platform's native controls and gestures. Custom chrome must earn its
place; a custom tab bar or modal is a product decision, not a default.

## 18. Indonesian by default

**Evidence.** The reference is English and US-local. Its structure transfers; its language
and local assumptions do not.

**Rule.** Every label, empty state, error, date and status is Indonesian. Use `WIB` for
event times, `2 jam lalu` for relative time, and Indonesian mosque/community imagery. Do
not translate Nextdoor copy literally.

## 19. Feedback at the point of action

**Evidence.** Nextdoor confirms success inline: a replaced control (“Message sent / See
conversation”), a banner over the column (“Invite sent”), or a pill beside the title
(“Saved”). It avoids a global corner toast.

**Rule.** Put confirmation near the control that caused it. Use `role="status"` and keep the
message short. A destructive success needs an undo path or a clear next step.

## 20. Nudges, not gates

**Evidence.** Onboarding uses a dismissible checklist card and group setup uses an inline
carousel; after signup the product does not block the feed behind a wizard.

**Rule.** Guidance is inline and dismissible. Gate only the action that genuinely requires a
precondition (publishing before validation), and explain it at that action.

## 21. Borders over shadows

**Evidence.** Resting cards use a 1px hairline and no shadow; shadows appear only on floating
layers (menus, modals, chat dock).

**Rule.** Use canvas contrast and borders for structure. A shadow means the element floats.

## Review question

For every UI decision, ask:

> Does this help a reader understand who published a kabar, what it says, and what they can
> do next — without adding noise, false urgency or a second product?

If the answer is no, the decision does not belong in the public surface.
