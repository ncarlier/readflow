import React, { useCallback, useState } from 'react'

import Loader from '../../../components/Loader'
import Button from '../../../components/Button'
import Panel from '../../../components/Panel'
import Icon from '../../../components/Icon'

const formats = [
  { value: 'html', label: 'HTML page', icon: 'code' },
  { value: 'offline', label: 'HTML page with images', icon: 'image' },
]

interface Props {
  download: (format: string) => Promise<any>
  onCancel: (e: any) => void
}

export default ({ download, onCancel }: Props) => {
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

  return (
    <Panel>
      {loading && <Loader blur />}
      <header>
        <h1>Download article as ...</h1>
      </header>
      <section>
        {formats.map((format) => (
          <Button
            key={`format_${format.value}`}
            variant="flat"
            style={{ width: '10rem', padding: '1rem' }}
            title={`Download article as ${format.label}`}
            onClick={() => handleDownloadArticle(format.value).then(onCancel)}
          >
            <Icon name={format.icon} />
            <p>{format.label}</p>
          </Button>
        ))}
      </section>
      <footer>
        <Button title="Cancel" onClick={onCancel}>
          Cancel
        </Button>
      </footer>
    </Panel>
  )
}
