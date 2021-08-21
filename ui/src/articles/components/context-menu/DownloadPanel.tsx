import React, { useCallback, useState } from 'react'

import LinkIcon from '../../../components/LinkIcon'
import Loader from '../../../components/Loader'

const formats = [
  { value: 'html', label: 'HTML file', icon: 'code' },
  { value: 'html-single', label: 'HTML file with images', icon: 'image' },
  { value: 'zip', label: 'ZIP file', icon: 'archive' },
  { value: 'epub', label: 'EPUB file', icon: 'menu_book' },
]

interface Props {
  only?: string
  download: (format: string) => Promise<any>
  onCancel: (e: any) => void
}

export default ({ download, onCancel, only }: Props) => {
  const [loading, setLoading] = useState(false)
  const handleDownloadArticle = useCallback(
    async (format: string) => {
      setLoading(true)
      try {
        await download(format)
      } finally {
        setLoading(false)
      }
    },
    [download]
  )

  const onlyFilter = (format: { value: string }) => (only ? format.value === only : true)

  if (loading) {
    return <Loader blur />
  }

  return (
    <ul>
      {formats.filter(onlyFilter).map((format) => (
        <li key={`format_${format.value}`}>
          <LinkIcon
            title="Download article as ..."
            icon={format.icon}
            onClick={() => handleDownloadArticle(format.value).then(onCancel)}
          >
            <span>{format.label}</span>
          </LinkIcon>
        </li>
      ))}
      <li>
        <LinkIcon title="Cancel" icon="cancel" onClick={onCancel}>
          <span>Cancel</span>
        </LinkIcon>
      </li>
    </ul>
  )
}
