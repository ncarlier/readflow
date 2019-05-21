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
  let dragStartX = 0
  let left = 0
  let dragged = false

  // FPS Limit
  let startTime = 0
  const fpsInterval = 1000 / 60

  // Refs
  const wrapperRef = useRef<HTMLDivElement>(null)
  const bgRef = useRef<HTMLDivElement>(null)
  const elementRef = useRef<HTMLDivElement>(null)

  const updatePosition = () => {
    if (dragged) requestAnimationFrame(updatePosition)
    const now = Date.now()
    const elapsed = now - startTime

    const $bg = bgRef.current
    const $el = elementRef.current
    if (dragged && elapsed > fpsInterval && $bg && $el) {
      $el.style.transform = `translateX(${left}px)`
      const opacity = Math.abs(left) / 100
      if (opacity < 1 && opacity.toFixed(2) !== $bg.style.opacity) {
        $bg.style.opacity = opacity.toFixed(2)
      }
      if (opacity >= 1) {
        $bg.style.opacity = '1'
      }
      startTime = Date.now()
    }
  }

  const onDragStart = (clientX: number) => {
    const $el = elementRef.current
    if ($el) {
      dragged = true
      dragStartX = clientX
      $el.className = styles.item
      startTime = Date.now()
      requestAnimationFrame(updatePosition)
    }
  }

  const onDragEnd = () => {
    const $el = elementRef.current
    const $wrapper = wrapperRef.current
    if (dragged && $el && $wrapper) {
      dragged = false
      if (left < $el.offsetWidth * threshold * -1) {
        left = -$el.offsetWidth * 2
        $wrapper.style.maxHeight = '0'
        onSwipe()
      } else {
        left = 0
      }

      $el.className = styles.bouncing
      $el.style.transform = `translateX(${left}px)`
    }
  }

  const onMouseMove = (evt: MouseEvent) => {
    const l = evt.clientX - dragStartX
    if (l < 0) {
      left = l
    }
  }

  const onTouchMove = (evt: TouchEvent) => {
    const touch = evt.targetTouches[0]
    const l = touch.clientX - dragStartX
    if (l < 0) {
      left = l
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
