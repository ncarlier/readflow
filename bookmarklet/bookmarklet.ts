export {}

declare global {
  interface Window {
    readflowBookmarklet: ReadflowBookmarklet
  }
}

const trapMouseEvent = function (callback: (el: HTMLElement) => void) {
  return function (ev: MouseEvent) {
    ev.cancelBubble = true
    ev.stopPropagation()
    const el = this as HTMLElement
    callback(el)
  }
}

const setRemoveStyle = function (el: HTMLElement) {
  el.style.background = 'orange'
  el.style.border = '2px dashed red'
  el.style.cursor = 'crosshair'
}

const setKeepStyle = function (el: HTMLElement) {
  el.style.background = 'greenyellow'
  el.style.border = '2px dashed green'
  el.style.cursor = 'crosshair'
}

const unsetMouseOverStyle = function (el: HTMLElement) {
  el.style.background = 'initial'
  el.style.border = 'initial'
  el.style.cursor = 'initial'
}

type HistoryItem = {
  parent: Node
  element: HTMLElement | ChildNode
  placeholder: Text
}

class ReadflowBookmarklet {
  private doc: Document
  private history: HistoryItem[]
  private alt: boolean

  constructor() {
    this.doc = document
    this.history = []
    this.alt = false
  }

  private clickElement(el: HTMLElement) {
    unsetMouseOverStyle(el)
    if (this.alt) {
      this.doc.body.childNodes.forEach((node) => {
        if (node.nodeType === Node.ELEMENT_NODE && node.nodeName !== 'SCRIPT') {
          const placeholder = this.doc.createTextNode('')
          this.history.push({
            parent: this.doc.body,
            element: node,
            placeholder,
          })
          node.parentNode.replaceChild(placeholder, node)
        }
      })
      this.doc.body.appendChild(el)
    } else {
      const placeholder = this.doc.createTextNode('')
      this.history.push({
        parent: el.parentNode,
        element: el,
        placeholder,
      })
      el.parentNode.replaceChild(placeholder, el)
    }
  }

  private undo() {
    if (this.history.length) {
      const item = this.history.pop()
      item.parent.replaceChild(item.element, item.placeholder)
    }
  }

  private toggleMouseOverStyle(el: HTMLElement) {
    this.alt ? setKeepStyle(el) : setRemoveStyle(el)
  }

  private registerEventsListeners() {
    const onmouseover = trapMouseEvent(this.toggleMouseOverStyle.bind(this))
    const onmouseout = trapMouseEvent(unsetMouseOverStyle)
    const onclick = trapMouseEvent(this.clickElement.bind(this))
    this.doc.body.querySelectorAll<HTMLElement>('*').forEach((node) => {
      node.onmouseover = onmouseover
      node.onmouseout = onmouseout
      node.onclick = onclick
    })
  }

  private registerKeyboardShortcuts() {
    this.doc.body.addEventListener('keydown', (ev: KeyboardEvent) => {
      if (ev.ctrlKey && ev.code === 'Backspace') {
        this.undo()
      }
      if (ev.ctrlKey && ev.code === 'Backslash') {
        this.alt = !this.alt
      }
    })
  }

  boot() {
    console.log('running bookmarklet...')
    this.registerEventsListeners()
    this.registerKeyboardShortcuts()
  }
}

window.readflowBookmarklet = new ReadflowBookmarklet()
