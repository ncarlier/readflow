import { Location } from 'history'
import React, { useCallback, useContext } from 'react'
import { useMutation } from '@apollo/client'
import { useModal } from 'react-modal-hook'
import { Link, useHistory, useLocation } from 'react-router-dom'

import ConfirmDialog from '../../components/ConfirmDialog'
import DropdownMenu from '../../components/DropdownMenu'
import Kbd from '../../components/Kbd'
import LinkIcon from '../../components/LinkIcon'
import { LocalConfigurationContext, SortBy, SortOrder } from '../../context/LocalConfigurationContext'
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

function revertSortOrder(order: SortOrder | null) {
  return order === 'asc' ? 'desc' : 'asc'
}

function revertSortBy(by: SortBy | null) {
  return by === 'key' ? 'stars' : 'key'
}

function revertStatus(status: string | null) {
  return status === 'unread' ? 'read' : 'unread'
}

function getSortByMessage(req: GetArticlesRequest) {
  return req.sortBy === 'stars' ? 'Sort by date' : 'Sort by stars'
}

function getSortOrderMessage(req: GetArticlesRequest) {
  if (req.sortBy === 'stars') {
    return req.sortOrder === 'asc' ? 'More stars first' : 'Less stars first'
  }
  return req.sortOrder === 'asc' ? 'Recent articles first' : 'Older articles first'
}

function getLocationWithSortByParam(loc: Location, by: SortBy) {
  const params = new URLSearchParams(loc.search)
  params.set('by', by)
  return { ...loc, search: params.toString() }
}

function getLocationWithSortOrderParam(loc: Location, order: SortOrder) {
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

  const updateLocalConfigSortBy = useCallback(() => {
    const { sorting } = localConfiguration
    const by = revertSortBy(req.sortBy)
    const key = req.category ? `cat_${req.category}` : variant
    if (!(Object.prototype.hasOwnProperty.call(sorting, key) && sorting[key].by === by)) {
      sorting[key].by = by
      updateLocalConfiguration({ ...localConfiguration, sorting })
    }
  }, [req, variant, localConfiguration, updateLocalConfiguration])

  const updateLocalConfigSortOrder = useCallback(() => {
    const { sorting } = localConfiguration
    const order = revertSortOrder(req.sortOrder)
    const key = req.category ? `cat_${req.category}` : variant
    if (!(Object.prototype.hasOwnProperty.call(sorting, key) && sorting[key].order === order)) {
      sorting[key].order = order
      updateLocalConfiguration({ ...localConfiguration, sorting })
    }
  }, [req, variant, localConfiguration, updateLocalConfiguration])

  const toggleSortBy = useCallback(
    (event: KeyboardEvent) => {
      event.preventDefault()
      updateLocalConfigSortBy()
      push(getLocationWithSortByParam(loc, revertSortBy(req.sortBy)))
      return false
    },
    [loc, req, push, updateLocalConfigSortBy]
  )

  const toggleSortOrder = useCallback(
    (event: KeyboardEvent) => {
      event.preventDefault()
      updateLocalConfigSortOrder()
      push(getLocationWithSortOrderParam(loc, revertSortOrder(req.sortOrder)))
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
    <DropdownMenu title="Page options...">
      <ul>
        <li>
          <LinkIcon onClick={refresh} icon="refresh">
            <span>Refresh</span>
            <Kbd keys="shift+r" onKeypress={refresh} />
          </LinkIcon>
        </li>
        {variant === 'starred' && (
          <li>
            <LinkIcon
              as={Link}
              to={getLocationWithSortByParam(loc, revertSortBy(req.sortBy))}
              onClick={updateLocalConfigSortBy}
              icon="swap_horiz"
            >
              <span>{getSortByMessage(req)}</span>
              <Kbd keys="shift+b" onKeypress={toggleSortBy} />
            </LinkIcon>
          </li>
        )}
        <li>
          <LinkIcon
            as={Link}
            to={getLocationWithSortOrderParam(loc, revertSortOrder(req.sortOrder))}
            onClick={updateLocalConfigSortOrder}
            icon="sort"
          >
            <span>{getSortOrderMessage(req)}</span>
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
