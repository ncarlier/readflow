export default function(hljs) {
  const LINE_COMMENT = hljs.COMMENT('//', '$', {
    contains: [{
      begin: /\\\n/
    }]
  })

  return {
    name: 'evalfilter',
    keywords: {
      $pattern: hljs.UNDERSCORE_IDENT_RE,
      literal: 'true false nil',
      keyword: 'case default else for foreach function if in local return switch while',
      built_in: 'between day float hour int keys len lower match max min minute month now panic printf print reverse seconds sort split sprintf string time trim type upper weekday year setText setTitle setHTML setCategory markAsRead markAsToRead sendNotification triggerWebhook disableGlobalNotification',
    },
    contains: [
      hljs.C_NUMBER_MODE,
      hljs.APOS_STRING_MODE,
      hljs.QUOTE_STRING_MODE,
      LINE_COMMENT,
      hljs.C_BLOCK_COMMENT_MODE,
    ]
  }
}
