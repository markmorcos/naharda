# dashboard

## ADDED Requirements

### Requirement: The dashboard SHALL finish the live-number and stats experience with zero layout shift
On a live update the hero number MUST animate (count up) using tabular figures so no layout shift
occurs, honoring `prefers-reduced-motion`; the dashboard MUST surface public `/v1/stats` figures; and
the webfont swap MUST cause no layout shift (metric-matched fallback).

#### Scenario: Live number counts up
- **WHEN** a live update changes the hero value
- **THEN** it animates from old to new with no layout shift, or swaps instantly under reduced-motion

#### Scenario: Public stats are visible
- **WHEN** a visitor views the dashboard
- **THEN** a strip shows current `/v1/stats` figures (requests served, data points, last-updated)

#### Scenario: No layout shift from fonts
- **WHEN** the page loads and the webfont swaps in
- **THEN** Cumulative Layout Shift remains 0 (metric-matched fallback)
