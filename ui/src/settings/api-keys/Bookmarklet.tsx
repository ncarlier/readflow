import React from 'react'

import { API_BASE_URL } from '../../constants'

function createBookmarkletScript(token: string) {
  const { origin } = document.location
  const cred = btoa('api:' + token)
  return `javascript:(function(){
FP_URL='${API_BASE_URL}';
FP_CRED='${cred}';
var js=document.body.appendChild(document.createElement('script'));
js.onerror=function(){alert('Sorry, unable to load bookmarklet.')};
js.src='${origin}/bookmarklet.js'})();
`
}

function createBookmarkletMarkup(token: string) {
  const script = createBookmarkletScript(token)
  return {
    __html: `<a title="Bookmark me!" href="${script}" onClick="alert('Don\\'t click on me! But drag and drop me to your toolbar.'); return false;">
<i class="material-icons">bookmark</i>
</a>
`
  }
}

interface Props {
  token: string
}

export default ({ token }: Props) => <div dangerouslySetInnerHTML={createBookmarkletMarkup(token)} />
