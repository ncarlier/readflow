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

export const getStaticProps = async ({ locale }) => {
  const { default: content } = await import(`../policies/terms_${locale}.md`)
  return { props: { content } }
}

export default TermsSimple
