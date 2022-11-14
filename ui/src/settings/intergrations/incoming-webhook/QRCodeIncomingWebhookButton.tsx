import React, { useCallback, useState } from 'react'
import { useModal } from 'react-modal-hook'

import { ButtonIcon, ConfirmDialog, Loader } from '../../../components'
import { fetchAPI, withCredentials } from '../../../helpers'
import { useMessage } from '../../../contexts'
import { useAuth } from '../../../auth/AuthProvider'

interface Props {
  token: string
}

export default ({ token }: Props) => {
  const { user } = useAuth()
  const [loading, setLoading] = useState(false)
  const [dataURI, setDataURI] = useState('')
  const { showErrorMessage } = useMessage()

  const [showQRCodeModal, hideQRCodeModal] = useModal(
    () => (
      <ConfirmDialog title="Incoming Webhook" onConfirm={hideQRCodeModal}>
        Scan this QR code with another device to allow someone else to send you articles from his readflow.
        <div style={{ textAlign: 'center' }}>
          <img src={dataURI} alt="QR code" />
        </div>
      </ConfirmDialog>
    ),
    [dataURI]
  )

  const generateQRCode = useCallback(async () => {
    setLoading(true)
    try {
      const headers = withCredentials(user)
      const res = await fetchAPI('/qr', { t: token }, { method: 'GET', headers })
      if (res.ok) {
        const data = await res.blob()
        setDataURI(window.URL.createObjectURL(data))
        showQRCodeModal()
      } else {
        const err = await res.json()
        throw new Error(err.detail || res.statusText)
      }
    } catch (err: any) {
      showErrorMessage(err.message)
    } finally {
      setLoading(false)
    }
  }, [token, showErrorMessage, showQRCodeModal])

  const handleClick = useCallback(() => {
    if (dataURI === '') {
      generateQRCode()
    } else {
      showQRCodeModal()
    }
  }, [dataURI, generateQRCode, showQRCodeModal])

  if (loading) {
    return <Loader center />
  }

  return <ButtonIcon title="Generate QR code" icon="qr_code" onClick={handleClick} />
}
