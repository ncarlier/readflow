import React, { SyntheticEvent, useCallback } from 'react'
import { Location } from 'history'
import { useMutation } from '@apollo/client'
import { useModal } from 'react-modal-hook'
import { Link, useHistory, useLocation } from 'react-router-dom'

import { ConfirmDialog, DropDownMenu, Kbd, LinkIcon } from '../../components'
import {
  DisplayMode,
  DisplayPreference,
  DisplayPreferences,
  SortBy,
  SortOrder,
  useLocalConfiguration,
} from '../../contexts/LocalConfigurationContext'
import { useMessage } from '../../contexts'
import { getGQLError } from '../../helpers'
import { GetArticlesRequest, MarkAllArticlesAsReadRequest, MarkAllArticlesAsReadResponse } from '../models'
import { MarkAllArticlesAsRead } from '../queries'
import { updateCacheAfterMarkAllAsRead } from '../cache'

type Variant = keyof DisplayPreferences

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

function revertDisplayMode(mode: DisplayMode | null) {
  return mode === 'grid' ? 'list' : 'grid'
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

function getDisplayModeMessage(mode: DisplayMode) {
  return mode === 'grid' ? 'Display as list' : 'Display as grid'
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

export const ArticlesPageMenu = (props: Props) => {
  const { refresh, req, variant } = props

  const loc = useLocation()
  const { push } = useHistory()

  const { showErrorMessage } = useMessage()
  const { localConfiguration, updateLocalConfiguration } = useLocalConfiguration()
  const [markAllArticlesAsReadMutation] = useMutation<MarkAllArticlesAsReadResponse, MarkAllArticlesAsReadRequest>(
    MarkAllArticlesAsRead
  )

  const getDisplayConfigKey = useCallback(() => {
    return req.category ? `cat_${req.category}` : variant
  }, [variant, req])

  const getDisplayPreference = useCallback((): DisplayPreference => {
    const { display } = localConfiguration
    const key = getDisplayConfigKey()
    if (Object.prototype.hasOwnProperty.call(display, key)) {
      return display[key]
    }
    return { by: 'key', order: 'asc', mode: 'list' }
  }, [localConfiguration, getDisplayConfigKey])

  const setDisplayPreference = useCallback(
    (pref: Partial<DisplayPreference>) => {
      const update = { ...getDisplayPreference(), ...pref }
      const { display } = localConfiguration
      display[getDisplayConfigKey()] = update
      updateLocalConfiguration({ ...localConfiguration, display })
    },
    [getDisplayPreference, localConfiguration, getDisplayConfigKey, updateLocalConfiguration]
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

  const toggleSortBy = useCallback(
    (event: SyntheticEvent | KeyboardEvent) => {
      event.preventDefault()
      const by = revertSortBy(req.sortBy)
      setDisplayPreference({ by })
      push(getLocationWithSortByParam(loc, revertSortBy(req.sortBy)))
    },
    [loc, req, push, setDisplayPreference]
  )

  const toggleSortOrder = useCallback(
    (event: SyntheticEvent | KeyboardEvent) => {
      event.preventDefault()
      const order = revertSortOrder(req.sortOrder)
      setDisplayPreference({ order })
      push(getLocationWithSortOrderParam(loc, revertSortOrder(req.sortOrder)))
    },
    [loc, req, push, setDisplayPreference]
  )

  const toggleDisplayMode = useCallback(
    (event: SyntheticEvent | KeyboardEvent) => {
      event.preventDefault()
      const mode = revertDisplayMode(getDisplayPreference().mode)
      setDisplayPreference({ mode })
    },
    [getDisplayPreference, setDisplayPreference]
  )

  const toggleStatus = useCallback(
    (event: SyntheticEvent | KeyboardEvent) => {
      event.preventDefault()
      push(getLocationWithStatusParam(loc, revertStatus(req.status)))
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
    <DropDownMenu title="Page options...">
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
              onClick={toggleSortBy}
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
            onClick={toggleSortOrder}
            icon="low_priority"
          >
            <span>{getSortOrderMessage(req)}</span>
            <Kbd keys="shift+o" onKeypress={toggleSortOrder} />
          </LinkIcon>
        </li>
        <li>
          <LinkIcon onClick={toggleDisplayMode} icon="dashboard">
            <span>{getDisplayModeMessage(getDisplayPreference().mode)}</span>
            <Kbd keys="shift+d" onKeypress={toggleDisplayMode} />
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
    </DropDownMenu>
  )
}
