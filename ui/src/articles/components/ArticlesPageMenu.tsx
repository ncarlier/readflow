import { Location } from 'history'
import React, { useCallback, useContext } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useModal } from 'react-modal-hook'
import { Link, useHistory, useLocation } from 'react-router-dom'

import ConfirmDialog from '../../components/ConfirmDialog'
import DropdownMenu from '../../components/DropdownMenu'
import Kbd from '../../components/Kbd'
import LinkIcon from '../../components/LinkIcon'
import { LocalConfigurationContext } from '../../context/LocalConfigurationContext'
import { MessageContext } from '../../context/MessageContext'
import { getGQLError } from '../../helpers'
import { GetArticlesRequest, MarkAllArticlesAsReadRequest, MarkAllArticlesAsReadResponse } from '../models'
import { MarkAllArticlesAsRead } from '../queries'
import { updateCacheAfterMarkAllAsRead } from '../cache'

type Variant = 'unread' | 'history' | 'starred' | 'offline'

interface Props {
  refresh: () => void
  req: GetArticlesRequest
  variant: Variant
}

function revertSortOrder(order: string | null) {
  return order === 'asc' ? 'desc' : 'asc'
}

function revertStatus(status: string | null) {
  return status === 'unread' ? 'read' : 'unread'
}

function getLocationWithSortParam(loc: Location, order: 'asc' | 'desc') {
  const params = new URLSearchParams(loc.search)
  params.set('sort', order)
  return { ...loc, search: params.toString() }
}

function getLocationWithStatusParam(loc: Location, status: 'read' | 'unread') {
  const params = new URLSearchParams(loc.search)
  params.set('status', status)
  return { ...loc, search: params.toString() }
}

export default (props: Props) => {
  const { refresh, req, variant } = props

  const loc = useLocation()
  const { push } = useHistory()

  const { showErrorMessage } = useContext(MessageContext)
  const { localConfiguration, updateLocalConfiguration } = useContext(LocalConfigurationContext)
  const [markAllArticlesAsReadMutation] = useMutation<MarkAllArticlesAsReadResponse, MarkAllArticlesAsReadRequest>(
    MarkAllArticlesAsRead
  )

  const markAllAsRead = useCallback(async () => {
    try {
      await markAllArticlesAsReadMutation({
        variables: { category: req.category },
        update: updateCacheAfterMarkAllAsRead,
      })
      await refresh()
    } catch (err) {
      showErrorMessage(getGQLError(err))
    }
  }, [markAllArticlesAsReadMutation, req, refresh, showErrorMessage])

  const updateLocalConfigSortOrder = useCallback(() => {
    const orders = localConfiguration.sortOrders
    const order = revertSortOrder(req.sortOrder)
    const key = req.category ? `cat_${req.category}` : variant
    if (!(Object.prototype.hasOwnProperty.call(orders, key) && orders[key] === order)) {
      orders[key] = order
      updateLocalConfiguration({ ...localConfiguration, sortOrders: orders })
    }
  }, [req, variant, localConfiguration, updateLocalConfiguration])

  const toggleSortOrder = useCallback(
    (event: KeyboardEvent) => {
      event.preventDefault()
      updateLocalConfigSortOrder()
      push(getLocationWithSortParam(loc, revertSortOrder(req.sortOrder)))
      return false
    },
    [loc, req, push, updateLocalConfigSortOrder]
  )

  const toggleStatus = useCallback(
    (event: KeyboardEvent) => {
      event.preventDefault()
      push(getLocationWithStatusParam(loc, revertStatus(req.status)))
      return false
    },
    [loc, req, push]
  )
  const [showMarkAllAsReadDialog, hideMarkAllAsReadDialog] = useModal(() => (
    <ConfirmDialog
      title="Mark all as read?"
      confirmLabel="Do it"
      onConfirm={() => markAllAsRead().then(hideMarkAllAsReadDialog)}
      onCancel={hideMarkAllAsReadDialog}
    >
      Are you sure to mark ALL articles as read?
    </ConfirmDialog>
  ))

  return (
    <DropdownMenu>
      <ul>
        <li>
          <LinkIcon onClick={refresh} icon="refresh">
            <span>Refresh</span>
            <Kbd keys="shift+r" onKeypress={refresh} />
          </LinkIcon>
        </li>
        <li>
          <LinkIcon
            as={Link}
            to={getLocationWithSortParam(loc, revertSortOrder(req.sortOrder))}
            onClick={updateLocalConfigSortOrder}
            icon="sort"
          >
            <span>{req.sortOrder === 'asc' ? 'Recent articles first' : 'Older articles first'}</span>
            <Kbd keys="shift+o" onKeypress={toggleSortOrder} />
          </LinkIcon>
        </li>
        {req.status === 'unread' && (
          <li>
            <LinkIcon onClick={showMarkAllAsReadDialog} icon="done_all">
              <span>Mark all as read</span>
              <Kbd keys="shift+m" onKeypress={showMarkAllAsReadDialog} />
            </LinkIcon>
          </li>
        )}
        {!!req.category && !!req.status && (
          <li>
            <LinkIcon as={Link} to={getLocationWithStatusParam(loc, revertStatus(req.status))} icon="history">
              <span>{req.status === 'read' ? 'View unread articles' : 'View read articles'}</span>
              <Kbd keys="shift+h" onKeypress={toggleStatus} />
            </LinkIcon>
          </li>
        )}
      </ul>
    </DropdownMenu>
  )
}
