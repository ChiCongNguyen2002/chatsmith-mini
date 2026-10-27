const { chromium } = require("playwright");

const BASE = process.env.BASE || "http://127.0.0.1:8080";

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage({ viewport: { width: 1280, height: 800 } });
  const errors = [];
  page.on("pageerror", (e) => errors.push("pageerror: " + e.message));
  page.on("console", (m) => { if (m.type() === "error") errors.push("console: " + m.text()); });

  await page.goto(BASE + "/", { waitUntil: "networkidle", timeout: 60000 });
  await page.click("#modelBtn");
  await page.waitForSelector("#picker .opt");
  await page.click("#modelBtn");

  await page.fill("#input", "Reply with one short sentence: what is 2+2?");
  await page.click("#send");
  await page.waitForFunction(() => /\$/.test(document.querySelector(".msg.assistant .chip")?.textContent || ""), null, { timeout: 90000 });

  console.log("reply:", await page.textContent(".msg.assistant .bubble"));
  console.log("chip:", await page.textContent(".msg.assistant .chip"));
  console.log("errors:", JSON.stringify(errors));

  await browser.close();
})().catch((e) => {
  console.error("ERROR:", e.message);
  process.exit(1);
});
