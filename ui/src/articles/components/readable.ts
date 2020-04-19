const css = `
:root {
  --text-color: #222;
  --bold-text-color: #111;
  --link-color: #009be5;

}
[data-theme="dark"] {
  --text-color: #ccc;
  --bold-text-color: #bbb;
  --link-color: #fff;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
  line-height: 1.6;
  color: var(--text-color);
  padding: 1rem 1rem 4rem 1rem;
  margin: auto;
}
img, figure {
  height: auto;
}
img, figure, iframe {
  display: block;
  max-width: 100%;
  margin: 0;
}
div {
  max-width: 100%;
}
pre, code {
  overflow: auto;
}
a {
  color: var(--link-color);
}
h1, h2, strong {
  color: var(--bold-text-color);
}
`

export default {
  css,
}
