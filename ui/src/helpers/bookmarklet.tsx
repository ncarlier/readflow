import { MouseEvent } from 'react'

import { API_BASE_URL } from '../constants'

export function getBookmarklet(apiKey: string) {
  const { origin } = document.location
  const cred = btoa('api:' + apiKey)
  return `javascript:(function(){
FP_URL="${API_BASE_URL}";
FP_CRED="${cred}";
var js=document.body.appendChild(document.createElement("script"));
js.onerror=function(){alert("Sorry, unable to load bookmarklet.")};
js.src="${origin}/bookmarklet.js"})();`
}

export function preventBookmarkletClick(e: MouseEvent<any>) {
  e.preventDefault()
  alert("Don't click on me! But drag and drop me to your toolbar.")
}
