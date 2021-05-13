import Layout from '@/components/Layout'
import Hero from '@/components/Hero'
import Features from '@/components/Features'
import Screenshot from '@/components/Screenshot'

const Home = () => (
  <Layout>
    <section>
      <Hero />
    </section>
    <section className="sea">
      <a id="features" />
      <Features />
    </section>
    <section>
      <a id="screenshot" />
      <Screenshot />
    </section>
  </Layout>
)

export default Home
