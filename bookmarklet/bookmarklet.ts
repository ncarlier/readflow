import * as styles from './styles'

export {}

declare global {
  interface Window {
    rfB: ReadflowBookmarklet
  }
}

type Article = {
  id?: number
  title?: string
  url?: string
  html?: string
  text?: string
  image?: string
}

type HistoryItem = {
  parent: Node | null
  element: HTMLElement | ChildNode
  placeholder: Text
}

const trapMouseEvent = function (callback: (el: HTMLElement, evt?: MouseEvent) => void) {
  return function (ev: MouseEvent) {
    ev.stopPropagation()
    ev.preventDefault()
    const el = this as HTMLElement
    callback(el, ev)
  }
}

const setMouserOverStyle = function (el: HTMLElement, evt?: MouseEvent) {
  Object.assign(el.style, evt && evt.ctrlKey ? styles.keep : styles.remove)
}

const unsetMouseOverStyle = function (el: HTMLElement) {
  Object.assign(el.style, styles.initial)
}

const addOpenGraphProps = function (doc: Document, article: Article): Article {
  return Array.from(doc.head.getElementsByTagName('meta')).reduce<Article>(
    (acc, $meta) => {
      const content = $meta.getAttribute('content')
      if (!content) {
        return acc
      }
      switch ($meta.getAttribute('property')) {
        case 'og:image':
          acc.image = content
          break
        case 'og:description':
          acc.text = content
          break
        case 'og:title':
          acc.title = content
          break
        default:
          break
      }
      return acc
    },
    article
  )
}

class ReadflowBookmarklet {
  private origin: string
  private doc: Document
  private history: HistoryItem[]
  private endpoint: string
  private key: string
  private popup: Window | null
  private controls: Node
  private article: Article | null

  constructor() {
    this.doc = document
    this.history = []
    this.article = null
  }

  private clickElement(el: HTMLElement, evt: MouseEvent) {
    unsetMouseOverStyle(el)
    if (evt.ctrlKey) {
      this.doc.body.childNodes.forEach((node) => {
        if (node.nodeType === Node.ELEMENT_NODE && node.nodeName !== 'SCRIPT' && node !== this.controls) {
          const placeholder = this.doc.createTextNode('')
          this.history.push({
            parent: this.doc.body,
            element: node,
            placeholder,
          })
          node.parentNode?.replaceChild(placeholder, node)
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
      el.parentNode?.replaceChild(placeholder, el)
    }
  }

  private undo() {
    if (this.history.length) {
      const item = this.history.pop()
      item?.parent?.replaceChild(item.element, item.placeholder)
    }
  }

  private registerEventsListeners() {
    const onmouseover = trapMouseEvent(setMouserOverStyle)
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
      if (ev.ctrlKey && ev.code === 'KeyZ') {
        this.undo()
      }
    })
  }

  private getContent() {
    let result = ''
    this.doc.body.childNodes.forEach((node) => {
      if (node.nodeType === Node.ELEMENT_NODE && node.nodeName !== 'SCRIPT' && node !== this.controls) {
        result += (node as HTMLElement).outerHTML
        // console.log(node)
      }
    })
    return result
  }

  private addControls() {
    // Build controls
    const controls = this.doc.createElement('div')
    Object.assign(controls.style, styles.controls)
    const frame = this.doc.createElement('iframe')
    frame.setAttribute('src', this.origin + '/bookmarklet.html')
    Object.assign(frame.style, styles.iframe)
    const drag = this.doc.createElement('div')
    Object.assign(drag.style, styles.drag)
    // Add drag support
    let dragging = false
    const offset = [0, 0]
    drag.onmousedown = (evt) => {
      dragging = true
      offset[0] = controls.offsetLeft - evt.clientX
      offset[1] = controls.offsetTop - evt.clientY
    }
    this.doc.onmouseup = () => (dragging = false)
    this.doc.onmousemove = (evt) => {
      evt.preventDefault()
      if (dragging) {
        controls.style.left = evt.clientX + offset[0] + 'px'
        controls.style.top = evt.clientY + offset[1] + 'px'
      }
    }
    // Add controls to the DOM
    controls.appendChild(frame)
    controls.appendChild(drag)
    this.controls = this.doc.body.appendChild(controls)
    this.popup = frame.contentWindow
  }

  private async post(article: Article) {
    const r = new Request(this.endpoint)
    const headers = new Headers({
      Accept: 'application/json',
      Authorization: 'Basic ' + this.key,
      'Content-Type': 'application/json',
    })
    const res = await fetch(r, {
      method: 'POST',
      headers,
      mode: 'cors',
      body: JSON.stringify(article),
    })
    if (res.ok) {
      console.debug('article added to readflow')
      return res.json()
    } else {
      throw `unable to send article: ${res.statusText}`
    }
  }

  private openResult() {
    if (this.article) {
      document.location.href = `${this.origin}/inbox/${this.article.id}`
    }
  }

  private onMessage(evt: MessageEvent) {
    const event = evt.data
    switch (event) {
      case 'content':
      case 'page':
        console.debug(`sending ${event}...`)
        this.popup?.postMessage('loading', '*')
        let article: Article = {
          title: document.title,
          url: document.location.href,
        }
        if (event === 'content') {
          article.html = this.getContent()
          article = addOpenGraphProps(this.doc, article)
        }
        this.post(article).then(
          (resp) => {
            this.article = resp.Articles[0]
            this.popup?.postMessage('success', '*')
          },
          (err) => {
            alert(err)
            this.popup?.postMessage('error', '*')
          }
        )
        break
      case 'openResult':
        this.openResult()
        break
      case 'close':
        this.close()
        break
      default:
        console.warn('unknown message event:', event)
    }
  }

  boot(origin: string, baseurl: string, key: string) {
    this.endpoint = baseurl + '/articles'
    this.origin = origin
    this.key = key
    this.registerEventsListeners()
    this.registerKeyboardShortcuts()
    this.addControls()
    window.addEventListener('message', this.onMessage.bind(this))
    console.debug('bookmarklet up and running...')
  }

  close() {
    this.doc.location.reload()
  }
}

window.rfB = new ReadflowBookmarklet()
