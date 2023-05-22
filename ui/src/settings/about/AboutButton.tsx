/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React from 'react'
import ReactModal from 'react-modal'
import { useModal } from 'react-modal-hook'

import { ButtonIcon, Logo, Panel } from '../../components'
import dialogStyles from '../../components/Dialog.module.css'
import version from '../../version'
import styles from './AboutButton.module.css'

interface Props {
  closeHandler: () => void
}

const AboutPanel = ({ closeHandler }: Props) => (
  <Panel className={styles.about}>
    <ButtonIcon title="close" onClick={closeHandler} icon="close" />
    <h1>
      <Logo name="readflow" />
    </h1>
    <span>({version})</span>
    <p>Read your Internet article flow in one place with complete peace of mind and freedom.</p>
    <ul>
      <li>
        <a href="https://github.com/ncarlier/readflow" rel="noreferrer noopener" target="_blank">
          Sources
        </a>
      </li>
      <li>
        <a href="https://github.com/ncarlier/readflow/issues" rel="noreferrer noopener" target="_blank">
          Bug or feature request
        </a>
      </li>
      <li>
        <a href="https://www.paypal.me/nunux" rel="noreferrer noopener" target="_blank">
          Support this project
        </a>
      </li>
    </ul>
  </Panel>
)

const AboutButton = () => {
  const [showAboutModal, hideAboutModal] = useModal(() => (
    <ReactModal
      isOpen
      shouldCloseOnEsc
      shouldCloseOnOverlayClick
      shouldFocusAfterRender
      onRequestClose={hideAboutModal}
      className={dialogStyles.dialog}
      overlayClassName={dialogStyles.overlay}
    >
      <AboutPanel closeHandler={hideAboutModal} />
    </ReactModal>
  ))

  return <ButtonIcon title="About" onClick={showAboutModal} icon="info" />
}

export default AboutButton
