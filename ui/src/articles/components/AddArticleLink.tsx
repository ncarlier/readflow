/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React, { SyntheticEvent, useCallback } from 'react'
import ReactModal from 'react-modal'
import { useModal } from 'react-modal-hook'
import { Link } from 'react-router-dom'

import { Category } from '../../categories/models'
import styles from '../../components/Dialog.module.css'
import { LinkIcon } from '../../components'
import { useKeyboard } from '../../hooks'
import { Article } from '../models'
import { AddArticleForm } from '.'

interface Props {
  category?: Category
  onSuccess: (article: Article) => void
}

export const AddArticleLink = ({ category, onSuccess }: Props) => {
  const [showAddModal, hideAddModal] = useModal(() => (
    <ReactModal
      isOpen
      shouldCloseOnEsc
      shouldCloseOnOverlayClick
      shouldFocusAfterRender
      onRequestClose={hideAddModal}
      className={styles.dialog}
      overlayClassName={styles.overlay}
      style={{ content: { minWidth: '50vw' } }}
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
    <LinkIcon icon="add" onClick={handleOnClick} as={Link} to={'/inbox/add'}>
      Add new ...
    </LinkIcon>
  )
}
