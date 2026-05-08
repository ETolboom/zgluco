# zgluco

`zgluco` exports your diabetes management data from Nightscout into clean, text-based tables optimized for LLM analysis.

## Rationale

LLMs excel at pattern recognition in natural language and structured text. When you hand them a graph or raw JSON, the results are mediocre. When you give them a well-formatted table of glucose values, treatments, and profile changes (with timestamps, units, and context) they can reason about your more data meaningfully.

`zgluco` bridges that gap: pull your data, get a clipboard-ready text export, paste it into your LLM of choice with a prompt.

## Getting Started

1. Download the executable for your platform (see Releases)
2. Run:
   ```
   zgluco format --all --days 7 --nightscout-url https://your-nightscout-instance.example.com --nightscout-api-key "apikey"
   ```
   Output is printed to stdout and copied to your clipboard automatically.
3. Paste the clipboard content into your LLM prompt

In case your Nightscout instance is public the API key can be omitted as well as other options:
```
zgluco format --sections sgv,treatments -n 5 --nightscout-url https://your-nightscout-instance.example.com
```

If you always use the same instance, you can drop the flags and set environmental variables or put your credentials in a `.env` file in the same directory instead:
```
NIGHTSCOUT_URL=https://your-nightscout-instance.example.com
NIGHTSCOUT_API_KEY=your_api_key
```

## What the Export Contains

The output is designed to give an LLM the right level of detail — not raw five-minute readings, but enough context to reason about patterns and cause and effect.

**Bucketed glucose.** Rather than exporting every CGM reading, SGVs are grouped into 15-minute buckets. Each row shows the average, min, and max glucose for that window plus the modal trend direction. This compresses a week of data to a manageable size while preserving the shape of excursions.

**Grouped treatments.** Consecutive zero temp basals are merged into a single entry with the combined duration, so a long low-glucose suspend doesn't flood the output with individual records. Consecutive SMBs within a 10-minute window are similarly aggregated into one line showing total insulin and the span they covered.

**Profile changelog.** The profile section lists your current basal rates, ISF, carb ratios, and glucose targets. If your profile has changed since the previous stored version, a changelog is appended showing exactly what was added, modified, or removed (and when). This way an LLM can distinguish a pattern caused by a bad profile setting from one that started after you changed it.

## Getting the Best Results

In this section I'll explain how I personally use LLMs (Mainly Claude) to analyze my glucose data. 

The heavy lifting is done by the [`closed_loop_analysis`](docs/closed_loop_analysis.md) skill. This is a large prompt that has instructions on how Claude should handle your data and what actionable insights it should provide.

It also covers some fundamental concepts from Loop, Trio (OpenAPS oref), and AndroidAPS, such that it tries to reason within each system's model. This way it won't apply Loop semantics to Trio data or vice versa.

You can add it by going to `Customize → Skills → Upload a Skill` and selecting the closed loop analysis Markdown file. Please note that you should download the raw file (Click the "download" or "raw" button) on the Markdown file as to include YAML frontmatter (contains the skill title and description)

If you are missing this you can always go `Customize → Skills → Write skill instructions` and copy-paste the prompt manually.

**Workflow:**

1. Run `zgluco format --all` and copy the output
2. Open a chat in a Claude **Project**, this provides a persistent context lets it reference previous recommendations across sessions
3. Add project instructions describing your setup: pump model, CGM sensor, insulin type. 
4. Paste the output of `zgluco` as a file and ask Claude to analyze it with the `closed-loop-analysis` skill. 

(The exports include context like large boluses during a high, profile changes, and temp targets, which lets the model reason about cause and effect rather than just raw numbers)

**NOTE: An LLM is not a medical professional. Please consider it a decision support. 
Always critically evaluate whether what it says 1) makes sense, 2) is reasonable 3) most importantly: is safe.**

**Tips:**

- Always try to mention any important details such as: "I was ill", "Last Sunday was a cheat day", "I am currently in a stressful period". 
By giving specific context you can avoid the LLM from hyperfocusing on that bad day and drawing wrong conclusions.
- Try to stick to one change at a time. The skill is created in such a way that it gives you recommendations based on its confidence and urgency. 
Tuning your profile is very much trial-and-error; don't change too many things at once, otherwise you won't know what actually what changes what.
- Export at least 5–7 days for overnight patterns, or 5–10 instances of a meal for meal-related tuning. It is difficult to give recommendations about say overnight basal
if you haven't had any meaningful data to find patterns in. 


## Supported Sources

| Source     | SGV | Treatments | Profile |
|------------|-----|------------|---------|
| Nightscout | Yes | Yes        | Yes     |

## Future Goals

- Add support for Tidepool
- Profile tuning based on ``oref0-autotune``

## Design Considerations

- Timestamps: all timestamps are assumed to be in UTC and are converted to the local machine's timezone.
- Units: all internal calculations use mg/dL (avoids precision loss; most hardware uses integers). Display respects your profile's preferred unit (mg/dL or mmol/L).
- Written in Go: compiles to a single native binary with no runtime dependencies. Handles large periods worth of CGM data in seconds without much of a struggle.

## Building from Source

Requires Go 1.21+.

```bash
go build ./cmd/zgluco
```

Or run directly:

```bash
go run ./cmd/zgluco format --all
```