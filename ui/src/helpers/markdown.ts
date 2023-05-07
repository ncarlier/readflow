export function basicMarkdownToHTML(text: string): string {
    // replace the linebreaks with <br>
    let result = text.replace(/(?:\r\n|\r|\n)/g, '<br>')
    // check for links [text](url)
    const links = result.match(/\[.*?\)/g)
    if (links != null && links.length > 0) {
        for (const link of links) {
            const label = link.match(/\[(.*?)\]/)
            const url = link.match(/\((.*?)\)/)
            if (label && url) {
                result = result.replace(link,`<a href="${url[1]}" target="_blank">${label[1]}</a>`)
            }
        }
    }
    return result
}
