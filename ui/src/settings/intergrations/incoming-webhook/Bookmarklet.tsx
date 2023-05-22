import React from 'react'

import { getAPIURL } from '../../../helpers'

function createBookmarkletScript(token: string) {
  const { origin } = document.location
  const cred = btoa('api:' + token)
  return `javascript:(function (d,u,c) {
  let done = false;
  const js = d.createElement('script');
  js.type = 'text/javascript';
  js.src = u;
  js.onerror = function() { alert('Sorry, unable to load bookmarklet.'); };
  js.onload = js.onreadystatechange = function () {
    if (!done && (!this.readyState || this.readyState === 'loaded' || this.readyState === 'complete')) {
      done = true;
      c();
    }
  };
  d.body.appendChild(js);
})(document, '${origin}/bookmarklet.js', function () {
  window.rfB.boot('${origin}', '${getAPIURL()}', '${cred}');
});
`
}

function createBookmarkletMarkup(token: string) {
  const script = createBookmarkletScript(token)
  return {
    __html: `<a title="Bookmark me!" href="${script}" onClick="alert('Don\\'t click on me! But drag and drop me to your toolbar.'); return false;">
<i class="material-icons">bookmark</i>
</a>
`,
  }
}

interface Props {
  token: string
}

export default ({ token }: Props) => <div dangerouslySetInnerHTML={createBookmarkletMarkup(token)} />
