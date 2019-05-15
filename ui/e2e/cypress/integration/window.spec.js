/// <reference types="Cypress" />

context('Window', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/login')
  })

  it('cy.title() - get the title', () => {
    cy.title().should('include', 'to read')
  })
})
