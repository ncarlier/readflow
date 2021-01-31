// Styles

type CssStyleObject = Partial<CSSStyleDeclaration>

export const iframe: CssStyleObject = {
  width: '100%',
  height: '100%',
  border: 'none',
  margin: '0',
}

export const controls: CssStyleObject = {
  position: 'fixed',
  top: '50px',
  right: '50px',
  width: '320px',
  height: '200px',
  zIndex: '999999999',
  boxShadow: 'gray 5px 5px 10px',
}

export const drag: CssStyleObject = {
  position: 'absolute',
  top: '0',
  left: '0',
  width: '90%',
  height: '2em',
  cursor: 'move',
}

export const remove: CssStyleObject = {
  background: 'orange',
  border: '2px dashed red',
  cursor: 'crosshair',
}

export const keep: CssStyleObject = {
  background: 'greenyellow',
  border: '2px dashed green',
  cursor: 'crosshair',
}

export const initial: CssStyleObject = {
  background: 'initial',
  border: 'initial',
  cursor: 'initial',
}
