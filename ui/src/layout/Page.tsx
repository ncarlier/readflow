import React, { ReactNode, FC, PropsWithChildren } from 'react'

import { classNames } from '../helpers'
import { usePageTitle } from '../hooks'
import { Appbar, Content } from '.'
import styles from './Page.module.css'

interface Props extends PropsWithChildren {
  title?: string
  subtitle?: string
  className?: string
  header?: ReactNode
  actions?: ReactNode
  scrollToTop?: boolean
}

export const Page: FC<Props> = (props) => {
  const { title, subtitle, className, children, actions, scrollToTop } = props
  const { header = <Appbar title={title}>{actions}</Appbar> } = props

  usePageTitle(title, subtitle)

  return (
    <section className={classNames(styles.page, className)}>
      <header>{header}</header>
      <Content scrollToTop={scrollToTop}>{children}</Content>
    </section>
  )
}
