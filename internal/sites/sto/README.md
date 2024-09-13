
```py
    def get_episodes(self, url:str):
        # /serie/stream/westworld/staffel-1/
        resp = self.s.get(self.base_url+url)
        resp.raise_for_status()

        soup = BeautifulSoup(resp.text, "html.parser")
        results = []
        for e in soup.find_all("a"):
            l = e.get("href")

            if l != None and str(l).startswith(url+"/episode-"):
                try:
                    results.index(l)
                except:
                    results.append(l)
        return results
    
   

    def get_seasons(self, url:str):
        # /serie/stream/westworld/
        resp = self.s.get(self.base_url+url)
        resp.raise_for_status()
        soup = BeautifulSoup(resp.text, "html.parser")
        results = []
        for e in soup.find_all("a"):
            l = e.get("href")
            if l != None and str(l).startswith(url+"/staffel-") and not str(l).__contains__("/episode-"):
                try:
                    results.index(l)
                except:
                    results.append(l)
        return results
```