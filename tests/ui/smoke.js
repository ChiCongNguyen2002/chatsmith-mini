const { chromium } = require("playwright");

const BASE = process.env.BASE || "http://127.0.0.1:8080";
const SESSION = process.env.SESSION_COOKIE || "";

(async () => {
  const browser = await chromium.launch();
  const context = await browser.newContext({ viewport: { width: 1280, height: 800 } });
  if (SESSION) {
    await context.addCookies([{ name: "chatsmith_session", value: SESSION, url: BASE }]);
  }
  const page = await context.newPage();
  const errors = [];
  page.on("pageerror", (e) => errors.push("pageerror: " + e.message));
  page.on("console", (m) => { if (m.type() === "error") errors.push("console: " + m.text()); });

  await page.goto(BASE + "/", { waitUntil: "networkidle", timeout: 60000 });

  if (!SESSION) {
    await page.waitForSelector('a[href="/auth/google/login"]', { timeout: 15000 });
    console.log("login screen: ok");
    console.log("errors:", JSON.stringify(errors));
    await browser.close();
    return;
  }

  await page.waitForSelector("#userName:not(:empty)", { timeout: 15000 });
  console.log("user:", await page.textContent("#userName"));
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
