# compliance Specification

## Purpose
GDPR-facing surfaces for the email-collecting dashboard. Created by archiving change
add-privacy-gdpr. Update Purpose after archive.
## Requirements
### Requirement: A privacy policy SHALL be published
A `/privacy` page (EN + ar-EG) MUST state the data collected (email, consent, locale, timestamp),
the lawful basis (consent), retention, the cookieless stance, and how to request deletion, and MUST
be linked from the footer and the signup consent.

#### Scenario: A visitor reads the policy
- **WHEN** a visitor opens `/privacy`
- **THEN** they see what is collected, why, how long it is kept, and how to delete it

### Requirement: Users SHALL be able to delete their data
The service MUST provide a way to delete a captured email on request (self-serve endpoint and a
contact address), removing it from storage.

#### Scenario: Deletion request
- **WHEN** a user requests deletion of their email
- **THEN** the email is removed from `signups` and the response does not reveal whether it existed
