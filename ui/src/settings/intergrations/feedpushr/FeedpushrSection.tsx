import React from 'react'

import { Button, Logo } from '../../../components'
import { useCurrentUser } from '../../../contexts'

const FeedpushrSection = () => {
  const user = useCurrentUser()
  if (!user || user.plan !== 'premium') {
    return null
  }

  return (
    <section>
      <header>
        <h2>
          <Logo name="feedpushr" style={{ maxWidth: '2em', verticalAlign: 'middle' }} />
          Feedpushr
        </h2>
        <Button as={'a'} href={`https://feedpushr.nunux.org/${user.hashid}`} target="_blank" title="Manage my feeds">
          Manage
        </Button>
      </header>
      <p>Thanks to your user plan you can manage RSS feeds with Feedpushr and send articles to your readflow.</p>
    </section>
  )
}

export default FeedpushrSection
