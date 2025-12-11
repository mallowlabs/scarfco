# scarfco

A Static Code Analysis tools Result Files COnverter.

This tool converts result files of static code analysis tools to the [Checkstyle](https://checkstyle.sourceforge.io/) format.
Currently supports above tools XMLs.

* [SpotBugs](https://spotbugs.github.io/)
* [FindBugs](http://findbugs.sourceforge.net/)
* [PMD](https://pmd.github.io/)
* [CPD](https://pmd.github.io/latest/pmd_userdocs_cpd.html)

## How to use

For example, you can use on GitHub Actions with Reveiewdog.


```yaml
name: CI

on: [pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    # ... (snipped)
    - name: Setup Reviewdog
      uses: reviewdog/action-setup@v1
    - name: Setup scarfco
      run: curl -sfL https://raw.githubusercontent.com/mallowlabs/scarfco/main/install.sh | sh -s
    - name: Run SpotBugs
      run: mvn spotbugs:spotbugs
    - name: Run Reviewdog (SpotBugs)
      env:
        REVIEWDOG_GITHUB_API_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      run: cat target/spotbugsXml.xml | ./bin/scanrfco | reviewdog -name=spotbugs -f=checkstyle -reporter=github-pr-review -diff="git diff ${{ github.event.pull_request.base.sha }}"
```

If you use PMD.

```yaml
    - name: Run PMD
      run: mvn pmd:pmd
    - name: Run Reviewdog (PMD)
      env:
        REVIEWDOG_GITHUB_API_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      run: cat target/pmd.xml | ./bin/scanrfco | reviewdog -name=pmd -f=checkstyle -reporter=github-pr-review -diff="git diff ${{ github.event.pull_request.base.sha }}"
```

If you use CPD.

```yaml
    - name: Run CPD
      run: mvn pmd:cpd
    - name: Run Reviewdog (CPD)
      env:
        REVIEWDOG_GITHUB_API_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      run: cat target/cpd.xml | ./bin/scanrfco | reviewdog -name=cpd -f=checkstyle -reporter=github-pr-review -diff="git diff ${{ github.event.pull_request.base.sha }}"
```

## How to build

Requirements:

* Go 1.24.0+

```shell
$ go build scarfco.go
```

## Author

* [@mallowlabs](https://github.com/mallowlabs)
