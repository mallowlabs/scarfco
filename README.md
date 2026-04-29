# scarfco

A Static Code Analysis tools Result Files COnverter.

This tool converts result files of static code analysis tools to [Checkstyle](https://checkstyle.sourceforge.io/) format or [SARIF 2.1.0](https://sarifweb.azurewebsites.net/) format.
Currently supports the following tools:

* [SpotBugs](https://spotbugs.github.io/)
* [FindBugs](http://findbugs.sourceforge.net/)
* [PMD](https://pmd.github.io/)
* [CPD](https://docs.pmd-code.org/latest/pmd_userdocs_cpd.html)
* [Checkstyle](https://checkstyle.sourceforge.io/)

## Usage

```
scarfco [-format <format>] [file]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-format` | `checkstyle` | Output format: `checkstyle` or `sarif` |

Input is read from stdin when no file argument is given.

## How to use

### With Reviewdog (Checkstyle format)

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
      uses: mallowlabs/scarfco@main
      with:
        version: latest
    - name: Run SpotBugs
      run: mvn spotbugs:spotbugs
    - name: Run Reviewdog (SpotBugs)
      env:
        REVIEWDOG_GITHUB_API_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      run: cat target/spotbugsXml.xml | scarfco | reviewdog -name=spotbugs -f=checkstyle -reporter=github-pr-review -diff="git diff ${{ github.event.pull_request.base.sha }}"
```

If you use PMD:

```yaml
    - name: Run PMD
      run: mvn pmd:pmd
    - name: Run Reviewdog (PMD)
      env:
        REVIEWDOG_GITHUB_API_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      run: cat target/pmd.xml | scarfco | reviewdog -name=pmd -f=checkstyle -reporter=github-pr-review -diff="git diff ${{ github.event.pull_request.base.sha }}"
```

If you use CPD:

```yaml
    - name: Run CPD
      run: mvn pmd:cpd
    - name: Run Reviewdog (CPD)
      env:
        REVIEWDOG_GITHUB_API_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      run: cat target/cpd.xml | scarfco | reviewdog -name=cpd -f=checkstyle -reporter=github-pr-review -diff="git diff ${{ github.event.pull_request.base.sha }}"
```

### With GitHub Code Scanning (SARIF format)

```yaml
name: CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      security-events: write
    steps:
    # ... (snipped)
    - name: Setup scarfco
      uses: mallowlabs/scarfco@main
      with:
        version: latest
    - name: Run SpotBugs
      run: mvn spotbugs:spotbugs
    - name: Convert to SARIF
      run: cat target/spotbugsXml.xml | scarfco -format sarif > spotbugs.sarif
    - name: Upload SARIF
      uses: github/codeql-action/upload-sarif@v4
      with:
        sarif_file: spotbugs.sarif
```

## How to build

Requirements:

* Go 1.24.0+

```shell
$ go build scarfco.go
```

## Author

* [@mallowlabs](https://github.com/mallowlabs)
