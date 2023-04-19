
export default function() {
  const FILTERS = ['urlquery', 'base64', 'json', 'html2text' ]
  const PROPS = ['id', 'title', 'text', 'html', 'url', 'image', 'href']
  const FILTER = {
    begin: /\|\s*\w+:?/,
    keywords: { name: FILTERS },
  }
  const PROP = {
    className: 'template-tag',
    beginKeywords: PROPS.join(' '),
    keywords: { name: PROPS },
    starts: {
      endsWithParent: true,
      contains: [ FILTER ],
      relevance: 0
    }
  }

  return {
    name: 'fasttemplate',
    aliases: [ 'fast' ],
    contains: [
      {
        className: 'template-variable',
        begin: /\{\{/,
        end: /\}\}/,
        contains: [
          PROP,
        ]
      }
    ]
  }
}
