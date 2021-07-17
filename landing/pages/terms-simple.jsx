import ReactMarkdown from 'react-markdown'

const TermsSimple = ({content}) => {
  return (
    <article>
      <ReactMarkdown>
        {content}
      </ReactMarkdown>
    </article>
  )
}

export async function getStaticProps(context) {
  const { locale } = context
  const content = await import(`../policies/terms_${locale}.md`)
  return {
    props: { content: content.default },
  }
}

export default TermsSimple
