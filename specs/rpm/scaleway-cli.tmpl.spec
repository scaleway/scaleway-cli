Name:           scw
Version:        %{version_}
Release:        1%{?dist}
Summary:        Scaleway CLI

License:        Apache License 2.0
URL:            https://github.com/scaleway/scaleway-cli
Source0:        https://github.com/scaleway/scaleway-cli/archive/refs/tags/v%{version}.tar.gz

%if 0%{?suse_version}
BuildRequires: go git
%else
BuildRequires:  golang git
%endif

Provides:       %{name} = %{version}

%description
Scaleway CLI

%global debug_package %{nil}

%prep
%autosetup -n scaleway-cli-%{version}

%define build_epoch %(git log -1 --pretty="%%ct")
%define build_date %(date --date="@%{build_epoch}" -u +"%%Y-%%m-%%dT%%H:%%M:%%SZ")
%define git_branch %(git symbolic-ref -q --short HEAD || echo HEAD)
%define git_commit %(git rev-parse --short HEAD)

%build
export CGO_ENABLED=0
LDFLAGS="-w -extldflags -static -X main.Version=%{version} -X main.BuildDate=%{build_date} -X main.GitBranch=%{git_branch} -X main.GitCommit=%{git_commit}"
GOPROXY=direct GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o %{name} cmd/scw/main.go


%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}


%files
%{_bindir}/%{name}


%changelog
