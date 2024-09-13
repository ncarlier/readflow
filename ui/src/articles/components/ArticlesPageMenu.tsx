import React, { useCallback } from 'react'
import { useMutation } from '@apollo/client'
import { useModal } from 'react-modal-hook'
import { useHistory, useLocation } from 'react-router-dom'

import { ConfirmDialog, DropDownMenu, Kbd, LinkIcon } from '../../components'
import {
  DisplayPreference,
  DisplayPreferences,
  SortBy,
  SortOrder,
  useLocalConfiguration,
} from '../../contexts/LocalConfigurationContext'
import { useMessage } from '../../contexts'
import { getGQLError } from '../../helpers'
import {
  ArticleStatus,
  GetArticlesRequest,
  MarkAllArticlesAsReadRequest,
  MarkAllArticlesAsReadResponse,
} from '../models'
import { MarkAllArticlesAsRead } from '../queries'
import { updateCacheAfterMarkAllAsRead } from '../cache'
import { DropDownMenuItem } from '../../components/DropDownMenuItem'
import { ToggleDisplayMode } from './ToggleDisplayMode'
import { ToggleSortOrder } from './ToggleSortOrder'
import { ToggleSortBy } from './ToggleSortBy'
import { ToggleView } from './ToggleView'

type Variant = keyof DisplayPreferences

interface Props {
  refresh: () => void
  req: GetArticlesRequest
  variant: Variant
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
        variables: { status: req.status || 'inbox', category: req.category },
        update: updateCacheAfterMarkAllAsRead,
      })
      await refresh()
    } catch (err) {
      showErrorMessage(getGQLError(err))
    }
  }, [markAllArticlesAsReadMutation, req, refresh, showErrorMessage])

  const handleChangeSortOrder = useCallback((order: SortOrder) => {
    setDisplayPreference({order})
    const params = new URLSearchParams(loc.search)
    params.set('sort', order)
    push({ ...loc, search: params.toString() })
  }, [loc, setDisplayPreference, push])
  
  const handleChangeSortBy = useCallback((by: SortBy) => {
    setDisplayPreference({by})
    const params = new URLSearchParams(loc.search)
    params.set('by', by)
    push({ ...loc, search: params.toString() })
  }, [loc, setDisplayPreference, push])
  
  const handleChangeView = useCallback((view: ArticleStatus | 'starred') => {
    const params = new URLSearchParams(loc.search)
    if (view == 'starred') {
      params.set('starred', 'true')
      params.delete('status')
    } else {
      params.set('status', view)
      params.delete('starred')
    }
      push({ ...loc, search: params.toString() })
  }, [loc, setDisplayPreference, push])

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
        {!!req.category && (
          <li>
            <DropDownMenuItem label='View'>
              <ToggleView value={req.starred ? 'starred' : req.status ?? 'inbox'} onChange={handleChangeView} kbs="shift+h" />
            </DropDownMenuItem>
          </li>
        )}
        <li>
          <DropDownMenuItem label='Display as'>
            <ToggleDisplayMode value={getDisplayPreference().mode} onChange={(mode) => setDisplayPreference({mode})} kbs="shift+d" />
          </DropDownMenuItem>
        </li>
        {variant === 'starred' && (
          <li>
            <DropDownMenuItem label='Sort by'>
              <ToggleSortBy value={getDisplayPreference().by} onChange={handleChangeSortBy} kbs="shift+b" />
            </DropDownMenuItem>
          </li>
        )}
        <li>
          <DropDownMenuItem label='Order '>
            <ToggleSortOrder value={getDisplayPreference().order} onChange={handleChangeSortOrder} kbs="shift+o" />
          </DropDownMenuItem>
        </li>
        {req.status === 'inbox' && (
          <li>
            <LinkIcon onClick={showMarkAllAsReadDialog} icon="done_all">
              <span>Mark all as read</span>
              <Kbd keys="shift+del" onKeypress={showMarkAllAsReadDialog} />
            </LinkIcon>
          </li>
        )}
      </ul>
    </DropDownMenu>
  )
}
