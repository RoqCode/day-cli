# `day` — Activity Logger: Todo-Liste

Ein passiver Activity-Logger in Go. Pings werden in einer SQLite-Datenbank gespeichert
und am Ende des Tages zu einem Report zusammengefasst (optional mit KI-Gruppierung).

---

## Phase 1: Projekt-Setup

- [x] Neues Git-Repo anlegen: `day`
- [x] `go mod init github.com/DEINNAME/day`
- [x] Verzeichnisstruktur anlegen:
  ```
  day/
    cmd/day/main.go          # Entrypoint
    internal/config/         # Config laden
    internal/db/             # SQLite Setup & Queries
    internal/ping/           # Ping-Logik
    internal/report/         # Report-Generierung
    internal/scope/          # fzf Scope-Picker
    internal/ai/             # KI-Wrapper
  ```
- [ ] Dependencies hinzufügen:
  - `modernc.org/sqlite` — pure Go SQLite, kein CGO
  - `github.com/BurntSushi/toml` — Config-Parser
  - `github.com/spf13/cobra` — CLI-Framework (Commands/Flags)
  - `github.com/sashabaranov/go-openai` — OpenAI Client

---

## Phase 2: Config (`internal/config/`)

- [ ] `~/.config/day/config.toml` definieren und parsen
- [ ] Struct für Config anlegen (`Config`, `AIConfig`, `ScopesConfig`)
- [ ] Default-Werte wenn Config nicht vorhanden

Beispiel `config.toml`:

```toml
[ai]
provider    = "openai"
model       = "gpt-4o-mini"
api_key_env = "OPENAI_API_KEY"  # Name der Env-Variable, nicht der Key selbst

[scopes]
predefined = ["daily", "refinement", "retro", "planning", "meeting", "NOTICKET"]

[day]
data_dir = "~/.local/share/day"  # optional override
```

---

## Phase 3: Datenbank (`internal/db/`)

- [x] SQLite-Verbindung aufbauen, `~/.local/share/day/day.db` erstellen
- [x] Schema via `CREATE TABLE IF NOT EXISTS` beim Start migrieren:
  ```sql
  CREATE TABLE IF NOT EXISTS pings (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    ts       DATETIME NOT NULL,
    activity TEXT NOT NULL,
    scope    TEXT,
    source   TEXT NOT NULL DEFAULT 'manual'
  );
  CREATE INDEX IF NOT EXISTS idx_pings_ts    ON pings(ts);
  CREATE INDEX IF NOT EXISTS idx_pings_scope ON pings(scope);
  ```
- [ ] Funktion: `InsertPing(ts time.Time, activity, scope, source string) error`
- [ ] Funktion: `GetPingsForDay(date time.Time) ([]Ping, error)`
- [ ] Funktion: `GetRecentScopes(n int) ([]string, error)` — für fzf-Vorschläge

---

## Phase 4: `day ping` (`internal/ping/` + `internal/scope/`)

- [ ] Scope-Picker mit fzf implementieren (Shell-out zu `fzf`):
  - Zuletzt verwendete Scopes oben (letzte 3 aus DB via `GetRecentScopes`)
  - Predefined Scopes aus Config darunter
  - `CUSTOM`-Option für freie Eingabe
- [ ] Zeitparser implementieren:
  - `--at 14:30` → HH:MM parsen, heutiges Datum verwenden
  - `--ago 30` → `time.Now().Add(-30 * time.Minute)`
  - Default: `time.Now()`
- [ ] `--silent` Flag: kein Output, kein fzf-Picker (für automatisierte Pings)
- [ ] `--scope` Flag: Scope direkt übergeben (für automatisierte Pings)
- [ ] Ping in DB schreiben

```bash
# Beispiel-Aufrufe
day ping "OAuth Flow weitergebaut"           # fzf-Picker erscheint
day ping "Stand-up" --at 09:15              # rückwirkend
day ping "Review" --ago 30                  # 30 min zurück
day ping --silent "gc: PROJ-42: fix" --scope "PROJ-42"  # automatisiert
```

---

## Phase 5: `day status`

- [ ] Anzahl Pings heute anzeigen
- [ ] Timestamp des ersten und letzten Pings heute
- [ ] Scope-Übersicht: wie viele Pings pro Scope heute

Beispiel-Output:

```
Today: 8 pings (09:04 – 16:45)

PROJ-42    ████████  5
daily             1
planning          1
meeting           1
```

---

## Phase 6: `day report` (plain)

- [ ] Alle Pings des Tages chronologisch ausgeben, gruppiert nach Scope
- [ ] Flag `--yesterday` unterstützen
- [ ] Argument `YYYY-MM-DD` für beliebigen Tag unterstützen
- [ ] Schönes Terminal-Format mit `fmt` (kein externes Framework nötig)

Beispiel-Output:

```
Friday, 21 February 2026 — 8 pings

PROJ-42
  09:04  git commit: add oauth flow
  10:31  OAuth redirect logic fertig
  11:15  git commit: fix token refresh

daily
  09:15  Stand-up

meeting
  14:00  Sync mit Team wegen Deployment
```

---

## Phase 7: `day report --ai`

- [ ] OpenAI-Client einbinden (`go-openai`)
- [ ] Pings als strukturierten Text-Prompt aufbereiten
- [ ] System-Prompt schreiben:
  > "Du bist ein Assistent der hilft, Arbeitszeiten zu dokumentieren.
  > Gruppiere die folgenden Aktivitäten nach Thema, schätze den ungefähren
  > Zeitaufwand pro Gruppe und formuliere kurze, prägnante Zusammenfassungen
  > die sich für einen Arbeitszeitbericht eignen."
- [ ] Antwort ausgeben
- [ ] Graceful error wenn kein API-Key gesetzt (hilfreiche Fehlermeldung)

---

## Phase 8: Installation

- [ ] `Makefile` mit `make install` → baut Binary nach `~/.local/bin/day`

```makefile
install:
	go build -o ~/.local/bin/day ./cmd/day
```

- [ ] Alternativ: `go install ./cmd/day` direkt nutzen

---

## Phase 9: Dotfiles-Integration

- [ ] **`common/.config/zsh/git/git-utils.zsh`** neu anlegen:
  - `_day_branch_scope()` — Branch-Präfix-Parser aus `gc.zsh` auslagern
    (erkennt `PROJ-42`, `NOTICKET`, `HOTFIX`, `BUGFIX`, leerer String als Fallback)

- [ ] **`gc.zsh`** anpassen:
  - `_day_branch_scope()` aus `git-utils.zsh` nutzen statt eigenem Parser
  - Nach erfolgreichem Commit silent pingen:
    ```zsh
    local _scope=$(_day_branch_scope)
    day ping --silent "gc: $final_msg" --scope "$_scope" &
    ```

- [ ] **`gs.zsh`** anpassen:
  - Nach Branch-Switch silent pingen:
    ```zsh
    local _scope=$(_day_branch_scope)
    day ping --silent "gs: switched to $target" --scope "$_scope" &
    ```

- [ ] **`gd.zsh`** anpassen:
  - Beim Start silent pingen:
    ```zsh
    local _scope=$(_day_branch_scope)
    day ping --silent "gd: diffview${range:+ $range}" --scope "$_scope" &
    ```

- [ ] **Neovim** (`init.lua` o.ä.):
  - `VimEnter` Autocmd mit `day ping --silent`:
    ```lua
    vim.api.nvim_create_autocmd("VimEnter", {
      callback = function()
        local file = vim.fn.expand("%:.")  -- relativer Pfad
        if file == "" then return end
        -- Scope aus Git-Branch ermitteln
        local branch = vim.fn.system("git symbolic-ref --quiet --short HEAD 2>/dev/null"):gsub("\n", "")
        local scope = branch:match("^[^/]+/([A-Za-z][A-Za-z0-9]+%-%d+)") or ""
        vim.fn.jobstart({ "day", "ping", "--silent", "nvim: " .. file, "--scope", scope })
      end
    })
    ```

---

## Bonus (optional)

- [ ] `day edit ID` — Ping nachträglich bearbeiten (scope, activity, ts)
- [ ] `day delete ID` — Ping löschen
- [ ] `day report --format json` — maschinenlesbarer Output
- [ ] Shell-Completion via `cobra` (`day completion zsh >> ~/.zshrc`)
- [ ] `day report --week` — Wochenübersicht

---

## Tipp für den Einstieg

Fang mit **Phase 3** (Datenbank) an — das ist das Herzstück.
Wenn `InsertPing` und `GetPingsForDay` funktionieren, macht alles andere
mehr Spaß weil du sofort echte Daten siehst.

Danach **Phase 4** (`day ping`) — sobald du Pings reinschreiben kannst,
ist das Tool bereits nutzbar.
