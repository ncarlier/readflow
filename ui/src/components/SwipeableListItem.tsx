import React, { FC, ReactNode, useCallback, useEffect, useRef } from 'react'
import { classNames } from '../helpers'

import styles from './SwipeableListItem.module.css'

interface Props {
  background?: ReactNode
  threshold?: number
  direction?: 'left' | 'right'
  onSwipe: () => void
}

const fpsInterval = 1000 / 60

export const SwipeableListItem: FC<Props> = ({
  children,
  background,
  onSwipe,
  direction = 'left',
  threshold = 0.3,
}) => {
  // Drag & Drop
  const dragStartX = useRef(0)
  const borderOffset = useRef(0)
  const dragged = useRef(false)

  // FPS Limit
  const startTime = useRef(0)

  // Refs
  const wrapperRef = useRef<HTMLDivElement>(null)
  const bgRef = useRef<HTMLDivElement>(null)
  const elementRef = useRef<HTMLDivElement>(null)

  const updatePosition = useCallback(() => {
    if (dragged.current) requestAnimationFrame(updatePosition)
    const now = Date.now()
    const elapsed = now - startTime.current

    const $bg = bgRef.current
    const $el = elementRef.current
    if (dragged.current && elapsed > fpsInterval && $bg && $el) {
      $el.style.transform = `translateX(${borderOffset.current}px)`
      const opacity = Math.abs(borderOffset.current) / 100
      if (opacity < 1 && opacity.toFixed(2) !== $bg.style.opacity) {
        $bg.style.opacity = opacity.toFixed(2)
      }
      if (opacity >= 1) {
        $bg.style.opacity = '1'
      }
      const trigger = direction == 'left' ? borderOffset.current : -borderOffset.current
      if (trigger < $el.offsetWidth * threshold * -1) {
        $bg.style.color = 'white'
      } else {
        $bg.style.color = 'rgba(255, 255, 255, 0.3)'
      }
      startTime.current = Date.now()
    }
  }, [dragged, startTime, bgRef, elementRef, borderOffset, threshold])

  const onDragStart = useCallback(
    (clientX: number) => {
      const $el = elementRef.current
      if ($el) {
        dragged.current = true
        dragStartX.current = clientX
        $el.className = styles.item
        startTime.current = Date.now()
        requestAnimationFrame(updatePosition)
      }
    },
    [elementRef, dragged, startTime, updatePosition]
  )

  const onDragEnd = useCallback(() => {
    const $el = elementRef.current
    const $wrapper = wrapperRef.current
    if (dragged.current && $el && $wrapper) {
      dragged.current = false
      const trigger = direction == 'left' ? borderOffset.current : -borderOffset.current
      if (trigger < $el.offsetWidth * threshold * -1) {
        const out = $el.offsetWidth * 2
        borderOffset.current = direction == 'left' ? -out : out
        $wrapper.style.maxHeight = '0'
        onSwipe()
      } else {
        borderOffset.current = 0
      }

      $el.className = styles.bouncing
      $el.style.transform = `translateX(${borderOffset.current}px)`
    }
  }, [elementRef, wrapperRef, dragged, borderOffset, onSwipe, threshold])

  const onMouseMove = useCallback(
    (evt: MouseEvent) => {
      const l = evt.clientX - dragStartX.current
      if ((direction == 'left' && l < 0) || (direction == 'right' && l > 0)) {
        borderOffset.current = l
      }
    },
    [borderOffset]
  )

  const onTouchMove = useCallback(
    (evt: TouchEvent) => {
      const touch = evt.targetTouches[0]
      const l = touch.clientX - dragStartX.current
      if ((direction == 'left' && l < 0) || (direction == 'right' && l > 0)) {
        borderOffset.current = l
      }
    },
    [borderOffset]
  )

  const onDragStartMouse = useCallback(
    (evt: React.MouseEvent) => {
      onDragStart(evt.clientX)
      const $el = wrapperRef.current
      if ($el) {
        $el.addEventListener('mousemove', onMouseMove)
      }
    },
    [onDragStart, wrapperRef, onMouseMove]
  )

  const onDragStartTouch = useCallback(
    (evt: React.TouchEvent) => {
      const touch = evt.targetTouches[0]
      onDragStart(touch.clientX)
      const $el = wrapperRef.current
      if ($el) {
        $el.addEventListener('touchmove', onTouchMove)
      }
    },
    [wrapperRef, onDragStart, onTouchMove]
  )

  const onDragEndMouse = useCallback(() => {
    const $el = wrapperRef.current
    if ($el) {
      $el.removeEventListener('mousemove', onMouseMove)
    }
    onDragEnd()
  }, [wrapperRef, onMouseMove, onDragEnd])

  const onDragEndTouch = useCallback(() => {
    const $el = wrapperRef.current
    if ($el) {
      $el.removeEventListener('touchmove', onTouchMove)
    }
    onDragEnd()
  }, [wrapperRef, onTouchMove, onDragEnd])

  useEffect(() => {
    window.addEventListener('mouseup', onDragEndMouse)
    window.addEventListener('touchend', onDragEndTouch)
    return () => {
      window.removeEventListener('mouseup', onDragEndMouse)
      window.removeEventListener('touchend', onDragEndTouch)
    }
  }, [onDragEndMouse, onDragEndTouch])

  return (
    <div className={styles.wrapper} ref={wrapperRef}>
      <div ref={bgRef} className={classNames(styles.background, styles[direction])}>
        {background ? background : <span>Action</span>}
      </div>
      <div ref={elementRef} onMouseDown={onDragStartMouse} onTouchStart={onDragStartTouch} className={styles.item}>
        {children}
      </div>
    </div>
  )
}
