/// <reference types="Cypress" />

context('Viewport', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/login')
  })

  it('set the viewport size and dimension', () => {
    cy.get('#navbar').should('be.visible')
    cy.get('#navbar-fog').should('not.be.visible')
    cy.viewport(480, 640)

    // the navbar should have collapse since our screen is smaller
    cy.get('#navbar-fog')
      .should('be.visible')
      .click({ force: true })
    cy.get('#navbar').should('not.be.visible')
    cy.get('#appbar-menu')
      .should('be.visible')
      .click()
    cy.get('#navbar').should('be.visible')
    cy.get('#navbar-fog').should('be.visible')
    cy.get('#navbar-link-settings').click()
    cy.get('#navbar').should('not.be.visible')
  })
})
