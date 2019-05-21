import React, { ReactNode, useEffect, useRef } from 'react'

import styles from './SwipeableListItem.module.css'

interface Props {
  children: ReactNode
  background?: ReactNode
  threshold?: number
  onSwipe: () => void
}

export default ({ children, background, onSwipe, threshold = 0.3 }: Props) => {
  // Drag & Drop
  const dragStartX = useRef(0)
  const left = useRef(0)
  const dragged = useRef(false)

  // FPS Limit
  const startTime = useRef(0)
  const fpsInterval = 1000 / 60

  // Refs
  const wrapperRef = useRef<HTMLDivElement>(null)
  const bgRef = useRef<HTMLDivElement>(null)
  const elementRef = useRef<HTMLDivElement>(null)

  const updatePosition = () => {
    if (dragged) requestAnimationFrame(updatePosition)
    const now = Date.now()
    const elapsed = now - startTime.current

    const $bg = bgRef.current
    const $el = elementRef.current
    if (dragged && elapsed > fpsInterval && $bg && $el) {
      $el.style.transform = `translateX(${left.current}px)`
      const opacity = Math.abs(left.current) / 100
      if (opacity < 1 && opacity.toFixed(2) !== $bg.style.opacity) {
        $bg.style.opacity = opacity.toFixed(2)
      }
      if (opacity >= 1) {
        $bg.style.opacity = '1'
      }
      startTime.current = Date.now()
    }
  }

  const onDragStart = (clientX: number) => {
    const $el = elementRef.current
    if ($el) {
      dragged.current = true
      dragStartX.current = clientX
      $el.className = styles.item
      startTime.current = Date.now()
      requestAnimationFrame(updatePosition)
    }
  }

  const onDragEnd = () => {
    const $el = elementRef.current
    const $wrapper = wrapperRef.current
    if (dragged && $el && $wrapper) {
      dragged.current = false
      if (left.current < $el.offsetWidth * threshold * -1) {
        left.current = -$el.offsetWidth * 2
        $wrapper.style.maxHeight = '0'
        onSwipe()
      } else {
        left.current = 0
      }

      $el.className = styles.bouncing
      $el.style.transform = `translateX(${left}px)`
    }
  }

  const onMouseMove = (evt: MouseEvent) => {
    const l = evt.clientX - dragStartX.current
    if (l < 0) {
      left.current = l
    }
  }

  const onTouchMove = (evt: TouchEvent) => {
    const touch = evt.targetTouches[0]
    const l = touch.clientX - dragStartX.current
    if (l < 0) {
      left.current = l
    }
  }

  const onDragStartMouse = (evt: React.MouseEvent) => {
    onDragStart(evt.clientX)
    window.addEventListener('mousemove', onMouseMove)
  }

  const onDragStartTouch = (evt: React.TouchEvent) => {
    const touch = evt.targetTouches[0]
    onDragStart(touch.clientX)
    window.addEventListener('touchmove', onTouchMove)
  }

  const onDragEndMouse = () => {
    window.removeEventListener('mousemove', onMouseMove)
    onDragEnd()
  }

  const onDragEndTouch = () => {
    window.removeEventListener('touchmove', onTouchMove)
    onDragEnd()
  }

  useEffect(() => {
    window.addEventListener('mouseup', onDragEndMouse)
    window.addEventListener('touchend', onDragEndTouch)
    return () => {
      window.removeEventListener('mouseup', onDragEndMouse)
      window.removeEventListener('touchend', onDragEndTouch)
    }
  }, [])

  return (
    <div className={styles.wrapper} ref={wrapperRef}>
      <div ref={bgRef} className={styles.background}>
        {background ? background : <span>Action</span>}
      </div>
      <div ref={elementRef} onMouseDown={onDragStartMouse} onTouchStart={onDragStartTouch} className={styles.item}>
        {children}
      </div>
    </div>
  )
}
