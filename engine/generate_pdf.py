"""
PDF Export Engine — System Design Interview Prep

Generic PDF generator that combines HTML content pages into a single A4 PDF.
Called by the Go dashboard via subprocess.

Usage:
    python3 engine/generate_pdf.py --config config.json

Config JSON format:
    {
        "title": "URL Shortener",
        "sections": [
            {"type": "problem", "title": "URL Shortener", "html_path": "...", "html_content": "..."},
            {"type": "appendix", "label": "A", "title": "Load Balancing", "html_path": "...", "html_content": "..."},
            ...
        ],
        "output": "output/url-shortener.pdf",
        "highlight_from": "url-shortener"  // optional: highlight context
    }

PDF Layout:
    - A4 portrait
    - Margin: left 2cm, bottom 0.5cm, top 0cm, right 0cm
    - Page numbers: bottom-right corner
    - Fonts: Inter (body), JetBrains Mono (code)
"""

import argparse
import json
import sys
from pathlib import Path


# CSS for PDF rendering
PDF_CSS = """
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap');

@page {
    size: A4 portrait;
    margin: 0cm 0cm 0.5cm 2cm;

    @bottom-right {
        content: counter(page);
        font-family: 'Inter', sans-serif;
        font-size: 9pt;
        color: #6B7280;
        padding-right: 1cm;
    }
}

body {
    font-family: 'Inter', sans-serif;
    font-size: 11pt;
    line-height: 1.7;
    color: #1F2937;
}

/* Appendix section break */
.appendix-break {
    page-break-before: always;
    padding-top: 1cm;
}

.appendix-header {
    font-size: 18pt;
    font-weight: 700;
    color: #7C3AED;
    margin-bottom: 0.5cm;
    padding-bottom: 0.3cm;
    border-bottom: 2px solid #7C3AED;
}

.appendix-label {
    font-size: 12pt;
    color: #9CA3AF;
    text-transform: uppercase;
    letter-spacing: 0.1em;
}

/* Include same macro styles as dashboard */
.say-box {
    padding: 10px 14px;
    margin: 12px 0;
    background: #ECFDF5;
    border-left: 3px solid #059669;
    border-radius: 0 8px 8px 0;
    font-size: 10pt;
}

.thought-cloud {
    padding: 10px 14px;
    margin: 12px 0;
    background: #F3F4F6;
    border-radius: 8px;
    font-size: 10pt;
    color: #6B7280;
}

.avoid-box {
    padding: 10px 14px;
    margin: 12px 0;
    background: #FEF2F2;
    border-left: 3px solid #DC2626;
    border-radius: 0 8px 8px 0;
    font-size: 10pt;
}

.key-takeaway {
    padding: 10px 14px;
    margin: 12px 0;
    background: #EFF6FF;
    border-left: 3px solid #2563EB;
    border-radius: 0 8px 8px 0;
    font-size: 10pt;
    font-weight: 600;
}

.phase-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin: 20px 0 12px;
    padding-bottom: 6px;
    border-bottom: 2px solid #6366F1;
}

.phase-number {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    background: #6366F1;
    color: white;
    font-weight: 700;
    font-size: 10pt;
}

.phase-title { font-size: 14pt; font-weight: 700; }
.phase-time { font-size: 8pt; color: #6366F1; }

.code-block {
    margin: 12px 0;
    background: #1E293B;
    border-radius: 6px;
    padding: 12px;
}

.code-block code {
    font-family: 'JetBrains Mono', monospace;
    font-size: 9pt;
    color: #E2E8F0;
    white-space: pre;
}

table {
    width: 100%;
    border-collapse: collapse;
    font-size: 9pt;
    margin: 12px 0;
}

th { background: #F9FAFB; font-weight: 600; }
th, td { padding: 6px 10px; border: 1px solid #E5E7EB; text-align: left; }

.compare-best { background: #ECFDF5; }
.compare-alt { background: #FFFBEB; }
.compare-nofit { background: #FEF2F2; }

.checklist { background: #ECFDF5; padding: 12px; border-radius: 8px; margin: 12px 0; }
.checklist li::before { content: "✓ "; color: #059669; font-weight: 700; }
.checklist ul { list-style: none; padding: 0; }

/* Contextual highlight in PDF */
.keyword-highlight {
    background: #FEF9C3 !important;
    border-left: 3px solid #F59E0B !important;
}
"""


def build_html(config: dict) -> str:
    """Build a single HTML document from config sections."""
    sections_html = []

    for section in config["sections"]:
        content = section.get("html_content", "")
        if not content and section.get("html_path"):
            path = Path(section["html_path"])
            if path.exists():
                content = path.read_text()

        if section["type"] == "problem":
            sections_html.append(f'<div class="problem-section">{content}</div>')
        elif section["type"] == "appendix":
            label = section.get("label", "")
            title = section.get("title", "")
            sections_html.append(f'''
                <div class="appendix-break">
                    <div class="appendix-label">Appendix {label}</div>
                    <div class="appendix-header">{title}</div>
                    {content}
                </div>
            ''')

    return f"""<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>{PDF_CSS}</style>
</head>
<body>
    <h1 style="font-size:24pt;margin-bottom:0.3cm;">{config['title']}</h1>
    {''.join(sections_html)}
</body>
</html>"""


async def generate_pdf(config: dict):
    """Generate PDF using Playwright."""
    from playwright.async_api import async_playwright

    html = build_html(config)
    output_path = config["output"]

    Path(output_path).parent.mkdir(parents=True, exist_ok=True)

    async with async_playwright() as p:
        browser = await p.chromium.launch()
        page = await browser.new_page()
        await page.set_content(html, wait_until="networkidle")
        await page.pdf(
            path=output_path,
            format="A4",
            margin={"top": "0cm", "right": "0cm", "bottom": "0.5cm", "left": "2cm"},
            print_background=True,
            display_header_footer=True,
            footer_template='<div style="font-size:9px;text-align:right;width:100%;padding-right:1cm;color:#6B7280;"><span class="pageNumber"></span></div>',
            header_template='<div></div>',
        )
        await browser.close()

    print(f"PDF generated: {output_path}")


def main():
    parser = argparse.ArgumentParser(description="Generate system design prep PDF")
    parser.add_argument("--config", required=True, help="Path to config JSON file")
    args = parser.parse_args()

    with open(args.config) as f:
        config = json.load(f)

    import asyncio
    asyncio.run(generate_pdf(config))


if __name__ == "__main__":
    main()
