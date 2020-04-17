/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React from 'react'
import ReactModal from 'react-modal'
import { useModal } from 'react-modal-hook'

import ButtonIcon from '../../components/ButtonIcon'
import dialogStyles from '../../components/Dialog.module.css'
import Panel from '../../components/Panel'
import { VERSION } from '../../constants'
import styles from './AboutButton.module.css'
import logo from './logo.svg'

interface Props {
  closeHandler: () => void
}

const AboutPanel = ({ closeHandler }: Props) => (
  <Panel className={styles.about}>
    <ButtonIcon title="close" onClick={closeHandler} icon="close" />
    <h1>
      <img src={logo} alt="logo" />
    </h1>
    <span>({VERSION})</span>
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

export default () => {
  const [showAboutModal, hideAboutModal] = useModal(() => (
    <ReactModal
      isOpen
      shouldCloseOnEsc
      shouldCloseOnOverlayClick
      shouldFocusAfterRender
      appElement={document.getElementById('root')!}
      onRequestClose={hideAboutModal}
      className={dialogStyles.dialog}
      overlayClassName={dialogStyles.overlay}
    >
      <AboutPanel closeHandler={hideAboutModal} />
    </ReactModal>
  ))

  return <ButtonIcon title="About" onClick={showAboutModal} icon="info" />
}
