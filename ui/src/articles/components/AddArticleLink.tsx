/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React, { SyntheticEvent, useCallback } from 'react'
import ReactModal from 'react-modal'
import { useModal } from 'react-modal-hook'
import { Link } from 'react-router-dom'

import { Category } from '../../categories/models'
import styles from '../../components/Dialog.module.css'
import LinkIcon from '../../components/LinkIcon'
import useKeyboard from '../../hooks/useKeyboard'
import { Article } from '../models'
import AddArticleForm from './AddArticleForm'

interface Props {
  category?: Category
  onSuccess: (article: Article) => void
}

export default ({ category, onSuccess }: Props) => {
  const [showAddModal, hideAddModal] = useModal(() => (
    <ReactModal
      isOpen
      shouldCloseOnEsc
      shouldCloseOnOverlayClick
      shouldFocusAfterRender
      appElement={document.getElementById('root')!}
      onRequestClose={hideAddModal}
      className={styles.dialog}
      overlayClassName={styles.overlay}
    >
      <AddArticleForm
        category={category}
        onSuccess={(a) => {
          onSuccess(a)
          hideAddModal()
        }}
        onCancel={hideAddModal}
      />
    </ReactModal>
  ))

  const handleOnClick = useCallback(
    (ev: SyntheticEvent) => {
      showAddModal()
      ev.preventDefault()
    },
    [showAddModal]
  )

  useKeyboard('+', showAddModal)

  return (
    <LinkIcon icon="add" onClick={handleOnClick} as={Link} to={'/unread/add'}>
      Add new article
    </LinkIcon>
  )
}
